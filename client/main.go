package client

import (
	"log/slog"
	"net/http"
	"net/url"

	"github.com/Jack-Gledhill/smtp2gotify/env"

	api "github.com/gotify/go-api-client/v2/client"
	"github.com/gotify/go-api-client/v2/gotify"
)

var (
	client *api.GotifyREST
	log    = slog.With("module", "http")
)

func init() {
	endpoint, err := url.Parse(env.Vars.Gotify.URL)
	if err != nil {
		panic("Invalid Gotify URL: " + err.Error())
	}

	client = gotify.NewClient(endpoint, &http.Client{})
	verRes, err := client.Version.GetVersion(nil)
	if err != nil {
		log.Error("Failed to fetch server version", err)
		return
	}

	log.Info("Connected to Gotify server", "version", verRes.Payload.Version)
}
