package server

import (
	"context"
	"log/slog"
	"time"

	"github.com/Jack-Gledhill/smtp2gotify/env"

	"github.com/emersion/go-smtp"
)

var (
	backend *Backend
	log     = slog.With("module", "smtp")
	server  *smtp.Server
)

func init() {
	// Prepare server
	backend = &Backend{}

	server = smtp.NewServer(backend)
	server.Addr = ":" + env.Vars.SMTP.Port
	server.Domain = env.Vars.SMTP.Host
	server.WriteTimeout = 10 * time.Second
	server.ReadTimeout = 10 * time.Second
	server.MaxMessageBytes = 1024 * 1024
	server.MaxRecipients = 50       // This is ignored (for now), but we'll set it high so messages aren't rejected
	server.AllowInsecureAuth = true // We're using PLAIN authentication, which is considered insecure
}

// Run starts the SMTP server and blocks the calling goroutine
func Run() error {
	// Gracefully shut down the server when the main thread terminates
	defer func() {
		if err := server.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()

	// Start the SMTP server
	log.Info("Starting server...")
	return server.ListenAndServe()
}
