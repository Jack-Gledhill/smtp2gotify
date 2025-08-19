package server

import (
	"github.com/emersion/go-smtp"
	"github.com/google/uuid"
)

type Backend struct{}

// NewSession is called after client greeting (EHLO, HELO).
func (b *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	id := uuid.NewString()
	log.Debug("Initialising new SMTP session", "id", id)

	return &Session{
		conn: c,
		log:  log.With("uuid", id),
		UUID: id,
	}, nil
}
