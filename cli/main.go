package main

import (
	"flag"
	"fmt"
	"net/smtp"
)

func main() {
	from := flag.String("from", "someone@example.com", "Sender email address")
	to := flag.String("to", "someone-else@example.com", "Recipient email address")
	host := flag.String("host", "localhost", "SMTP host address")
	port := flag.String("port", "25", "SMTP port")
	username := flag.String("username", "admin", "SMTP username")
	password := flag.String("password", "test123", "SMTP password")
	flag.Parse()

	auth := smtp.PlainAuth("", *username, *password, *host)
	msg := []byte("From: " + *from + "\r\n" +
		"To: " + *to + "\r\n" +
		"Subject: smtp2gotify Test Message\r\n" +
		"\r\n" +
		"This is a test message sent to smtp2gotify via SMTP and then relayed to Gotify. If you're seeing this, that means things are working!\r\n")

	err := smtp.SendMail(*host+":"+*port, auth, *from, []string{*to}, msg)
	if err != nil {
		panic(err)
	}
	fmt.Println("Message sent successfully!")
}
