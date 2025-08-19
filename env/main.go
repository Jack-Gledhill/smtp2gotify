package env

import "os"

// Vars is populated at runtime with all the environment variables required by this program
var Vars *Environment

// Environment wraps all env vars recognised by this program
type Environment struct {
	SMTP   SMTP
	Gotify Gotify
}

// SMTP holds information about the SMTP server
type SMTP struct {
	Host     string
	Port     string
	Username string
	Password []byte
}

// Gotify holds environment variables about the Gotify server
type Gotify struct {
	APIToken string
	URL      string
}

func init() {
	Vars = &Environment{
		SMTP: SMTP{
			Host:     os.Getenv("SMTP_HOST"),
			Port:     os.Getenv("SMTP_PORT"),
			Username: os.Getenv("SMTP_USERNAME"),
			Password: []byte(os.Getenv("SMTP_PASSWORD")),
		},
		Gotify: Gotify{
			APIToken: os.Getenv("GOTIFY_API_TOKEN"),
			URL:      os.Getenv("GOTIFY_URL"),
		},
	}
}
