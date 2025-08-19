package server

import (
	"io"
	"log/slog"
	"strings"

	"github.com/Jack-Gledhill/smtp2gotify/client"
	"github.com/Jack-Gledhill/smtp2gotify/env"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"golang.org/x/crypto/bcrypt"
)

// Session handles an active connection from a client
type Session struct {
	conn *smtp.Conn
	log  *slog.Logger

	// UUID is an ephemeral identifier for the session, useful for filtering logs
	UUID string

	// Authenticated will be false until the client has passed an authentication challenge
	// The client cannot send MAIL, RCPT, or DATA commands until this field is true
	Authenticated bool

	// Message details provided by the client during the session
	From       string
	Options    *smtp.MailOptions
	Recipients []Recipient
	Body       []byte
}

// AuthMechanisms returns an array of valid auth mechanisms
func (s *Session) AuthMechanisms() []string {
	return []string{sasl.Plain}
}

// Auth is the handler for supported authenticators
func (s *Session) Auth(_ string) (sasl.Server, error) {
	s.log.Debug("Received AUTH command")
	return sasl.NewPlainServer(s.PlainAuth), nil
}

// PlainAuth is the handler for the PLAIN authentication mechanism.
// It compares the provided username against the one configured, and the provided password against the configured bcrypt hash.
// Once authenticated, it sets Session.Authenticated to true, allowing further commands to progress.
func (s *Session) PlainAuth(_ string, u string, p string) error {
	err := bcrypt.CompareHashAndPassword(env.Vars.SMTP.Password, []byte(p))
	if err != nil || u != env.Vars.SMTP.Username {
		s.log.Debug("Client failed authentication", "username", u)
		return smtp.ErrAuthFailed
	}

	s.log.Info("Client authenticated", "username", u)
	s.Authenticated = true
	return nil
}

// Mail is called when a client wishes to send a new message. It contains the sender address.
func (s *Session) Mail(f string, o *smtp.MailOptions) error {
	s.log.Debug("Received MAIL command", "from", f)

	if !s.Authenticated {
		s.log.Warn("Client attempted to send MAIL without authentication")
		return smtp.ErrAuthRequired
	}

	s.From = f
	s.Options = o
	return nil
}

// Rcpt adds a new recipient to the current message. A client will call this once for every recipient.
func (s *Session) Rcpt(to string, o *smtp.RcptOptions) error {
	s.log.Debug("Received RCPT command", "to", to)

	if !s.Authenticated {
		s.log.Warn("Client attempted to send RCPT without authentication")
		return smtp.ErrAuthRequired
	}

	r := Recipient{
		Address: to,
		Options: o,
	}

	s.Recipients = append(s.Recipients, r)
	return nil
}

// Data marks the beginning of the message's body. This function will block until the client has finished sending the body.
func (s *Session) Data(r io.Reader) error {
	s.log.Debug("Received DATA command")

	if !s.Authenticated {
		s.log.Warn("Client attempted to send DATA without authentication")
		return smtp.ErrAuthRequired
	}

	b, err := io.ReadAll(r)
	if err != nil {
		s.log.Error("Error reading data", "error", err)
		return err
	}

	s.Body = b
	s.Flush()
	return nil
}

// Reset is called whenever a client wishes to abort.
// All current message information is discarded, allowing the client to try again.
// Any previous authentication is not cleared, so the client can continue to send messages without re-authenticating.
func (s *Session) Reset() {
	s.log.Debug("Received RSET command")

	s.From = ""
	s.Recipients = []Recipient{}
	s.Body = []byte{}
}

// Logout is sent by a client when it is ready to terminate the connection completely.
// If a client wishes to connect again, it must re-authenticate.
func (s *Session) Logout() error {
	s.log.Debug("Received QUIT command")
	s.Reset()
	return nil
}

// SplitBody splits the message body into two parts: headers and content.
// It does this by looking for the first double CRLF ("\r\n\r\n") in the body, and then splitting the content either side of it.
func (s *Session) SplitBody() (string, string) {
	parts := strings.Split(string(s.Body), "\r\n\r\n")
	return parts[0], parts[1]
}

// GetHeaders retrieves a map of headers from the message body.
// Headers are returned as key-value pairs, where the key is the header name and the value is the header value.
// This function assumes that headers are separated by newlines and that the end of headers is indicated by a blank line.
func (s *Session) GetHeaders() map[string]string {
	stringHeaders, _ := s.SplitBody()
	headers := make(map[string]string)

	lines := strings.Split(stringHeaders, "\n")
	for _, line := range lines {
		parts := strings.Split(line, ": ")
		headers[parts[0]] = parts[1]
	}

	return headers
}

// GetSubject retrieves the "Subject" header from the message.
// Under the hood, this calls GetHeaders() then retrieves the "Subject" key from the returned map.
// If the "Subject" header is not present, it returns an empty string.
func (s *Session) GetSubject() string {
	sub, ok := s.GetHeaders()["Subject"]
	if !ok {
		return ""
	}

	return sub
}

// GetContent extracts the content of the message body, excluding all headers and the header separator line.
// This assumes that all text after the first double CRLF ("\r\n\r\n") is the content of the message.
func (s *Session) GetContent() string {
	_, content := s.SplitBody()
	return content
}

// Flush sends the current message to the Gotify client for processing.
// This should be called after Data completes.
func (s *Session) Flush() {
	s.log.Info("Received new message", "from", s.From)
	client.SendMessage(s.GetSubject(), s.GetContent())
}
