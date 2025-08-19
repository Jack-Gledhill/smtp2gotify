package main

import (
	"log/slog"

	"github.com/Jack-Gledhill/smtp2gotify/server"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.Debug("Running in debug mode")

	// Start the SMTP server
	err := server.Run()
	if err != nil {
		panic(err)
	}
}
