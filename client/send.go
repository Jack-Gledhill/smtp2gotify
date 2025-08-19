package client

import (
	"github.com/Jack-Gledhill/smtp2gotify/env"

	"github.com/gotify/go-api-client/v2/auth"
	"github.com/gotify/go-api-client/v2/client/message"
	"github.com/gotify/go-api-client/v2/models"
)

// SendMessage accepts a title and message string, then sends it to the Gotify server.
func SendMessage(t string, m string) {
	params := message.NewCreateMessageParams()
	params.Body = &models.MessageExternal{
		Title:    t,
		Message:  m,
		Priority: 0,
	}

	_, err := client.Message.CreateMessage(params, auth.TokenAuth(env.Vars.Gotify.APIToken))
	if err != nil {
		log.Error("Failed to send message", "error", err)
		return
	}

	log.Debug("Sent message to Gotify server")
}
