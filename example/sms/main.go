package main

import (
	"context"
	"log"
	"os"

	"github.com/gurza/go-redsms/redsms"
)

const Sender string = "REDSMS"

func main() {
	login := os.Getenv("REDSMS_LOGIN")
	if login == "" {
		log.Fatal("Unauthorized: REDSMS_LOGIN env is not set")
	}
	apiKey := os.Getenv("REDSMS_APIKEY")
	if apiKey == "" {
		log.Fatal("Unauthorized: REDSMS_APIKEY env is not set")
	}
	tp := redsms.SimpleAuthTransport{
		Login:  login,
		APIKey: apiKey,
	}
	c := redsms.NewClient(tp.Client())
	msg := redsms.Message{
		From:  Sender,
		To:    "+1234567890",
		Text:  "Hello world!",
		Route: redsms.MessageRouteSMS,
	}
	report, _, err := c.Message.Send(context.Background(), &msg)
	if err != nil {
		log.Fatal(err)
	}

	if len(report.Items) == 0 {
		log.Printf("Message not sent to anyone")
	}
	for _, item := range report.Items {
		log.Printf("Message send ID to %s: %s", item.To, item.UUID)
	}
	for _, e := range report.Errors {
		log.Printf("Error sending message to %s: %s", e.To, e.Message)
	}
}
