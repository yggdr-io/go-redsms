// The balance command utilize go-redsms as a cli tool for
// sending sms message.
package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/yggdr-corp/go-redsms/v2/redsms"
)

const Sender string = "REDSMS"

var (
	to   = flag.String("to", "", "recipient's mobile phone")
	text = flag.String("text", "", "message text")
)

func main() {
	flag.Parse()
	login := os.Getenv("REDSMS_LOGIN")
	if login == "" {
		log.Fatal("Unauthorized: REDSMS_LOGIN env is not set")
	}
	apiKey := os.Getenv("REDSMS_APIKEY")
	if apiKey == "" {
		log.Fatal("Unauthorized: REDSMS_APIKEY env is not set")
	}
	if *to == "" {
		log.Fatal("No recipient: you must specify recipient's mobile phone")
	}
	if *text == "" {
		log.Fatal("No text: you must specify message text")
	}
	ctx := context.Background()
	tp := redsms.SimpleAuthTransport{
		Login:  login,
		APIKey: apiKey,
	}
	client := redsms.NewClient(tp.Client())

	msg := redsms.Message{
		From:  Sender,
		To:    *to,
		Text:  *text,
		Route: redsms.MessageRouteSMS,
	}
	report, _, err := client.Message.Send(ctx, &msg)
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
