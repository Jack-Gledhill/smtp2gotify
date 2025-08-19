package server

import "github.com/emersion/go-smtp"

// Recipient represents an email recipient in the SMTP session.
// This is currently unused.
type Recipient struct {
	Address string
	Options *smtp.RcptOptions
}
