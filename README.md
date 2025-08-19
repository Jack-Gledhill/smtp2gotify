# SMTP to Gotify Relay

This is a very simple SMTP to Gotify adapter written in Go. 
It is essentially just a barebones SMTP server that translates received emails into Gotify messages and then forwards them to a Gotify server.

## What is Gotify?

[Gotify](https://gotify.net) is a self-hosted iOS/Android notification server that allows you to send messages to your mobile devices. 
It's a more modern alternative to email notifications, and it has a nice web interface for managing your notifications. 

I prefer Gotify over email for notifications, because it reduces the noise in my inbox and allows me to manage notifications more effectively. 
The only caveat is that Gotify is less supported, which is what this project aims to fix.

## But y tho?

This relay fixes a problem I've had with [my homelab](https://github.com/constellation-net): various services (e.g. TrueNAS) support SMTP for notifications, but don't support Gotify. 
In an effort to unionise the way I receive administrative notifications, I wrote this relay to do the translation for me. 

## Installation

You can build this relay from source (if you really must) with Go >=v1.23.2, or you can use the pre-built Docker image (like a sane person). 
See [example.compose.yml](example.compose.yml) for an example of how to run smtp2gotify with Docker Compose.

## Configuration

Configuration is done via environment variables. They're explained in this table:

| Environment variable | Description                                                           | Example value                |
|----------------------|-----------------------------------------------------------------------|------------------------------|
| `SMTP_HOST`          | The hostname of the SMTP server                                       | `smtp2gotify.example.com`    |
| `SMTP_PORT`          | The port that the SMTP server will listen on                          | `25`                         |
| `SMTP_USERNAME`      | Username that clients should use to authenticate with the SMTP server | `gotify`                     |
| `SMTP_PASSWORD`      | Bcrypt hash of the password that SMTP clients should provide          | Bcrypt hash                  |
| `GOTIFY_URL`         | The URL of the Gotify server to send messages to                      | `https://gotify.example.com` |
| `GOTIFY_API_TOKEN`   | API token used to authenticate with the Gotify server                 | `abcdef1234567890`           |

To get a Gotify API token, you'll need to create a new application on your Gotify server's web interface, then copy the API token from the application settings.

### Calculating a bcrypt hash

If you're on Linux or macOS, you can use the `htpasswd` command to generate a bcrypt hash of your password:

```shell
htpasswd -bnBC 10 "" YOUR_PASSWORD | tr -d ':\n'
```

If you're on Windows, you'll just have to figure it out for yourself. 

### Configuring Senders

It should be pretty straightforward to configure your SMTP clients. 
The only thing you really need to do is make sure that your client is **not** using TLS or STARTTLS (unless you're behind a reverse proxy that handles TLS termination), and that it is using the PLAIN authentication mechanism.

## Testing

This repository includes an SMTP test client that you can use to check everything's working correctly. 
It's located at `cli/main.go`, and you can run it with the following command:

```shell
go run cli/main.go \
    -host=localhost \
    -port=25 \
    -username=admin \
    -password=password \
    -from=someone@example.com \
    -to=someone-else@example.com
```

All being well, you should see a message in the CLI that the message was sent, and you should also see a message appear in your Gotify server.

## Security Considerations

Please note that this relay is NOT designed to be exposed on the internet.
It uses PLAIN authentication which is not very secure. 
It also doesn't use TLS, which means that traffic is not encrypted in-flight.
You should also ensure your Gotify server is secured with an SSL certificate, especially if it's accessible over the internet. 

You should only use this relay in an environment where you trust the senders.
While the relay does have a basic authentication mechanism, it is **not** designed for situations where sharing the SMTP credentials is a risk. 
If you don't feel comfortable sharing a single set of SMTP credentials all of your senders, then don't use this relay. 

This relay does not perform any spam filtering, validation or modification of the messages it receives, with the sole exception of translating the email into a Gotify message. 

## Dependencies

I had hoped to make this a fully dependency-less project, using only Go's standard library. 
Alas, `net/smtp` only has an SMTP client, not a server, and is permanently in a frozen state. 
Plus, with dependencies we can use bcrypt for hashing the SMTP password rather than storing it in plaintext, along with Gotify's official Go client library. 

The full list of direct dependencies is as follows:

- [github.com/emersion/go-smtp](https://github.com/emersion/go-smtp): SMTP server implementation
- [github.com/emersion/go-sasl](https://github.com/emersion/go-sasl): authentication mechanism for the SMTP server
- [golang.org/x/crypto/bcrypt](https://golang.org/x/crypto/bcrypt): validating bcrypt hashes
- [github.com/gotify/go-api-client](https://github.com/gotify/go-api-client): HTTP client for Gotify
- [github.com/google/uuid](https://github.com/google/uuid): issues UUIDs for each SMTP session

See `go.mod` for the list of transitive dependencies.