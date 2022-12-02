// The balance command utilize go-redsms as a cli tool for
// getting balance information.
package main

import (
	"context"
	"log"
	"os"

	"github.com/yggdr-corp/go-redsms/v1/redsms"
)

func main() {
	login := os.Getenv("REDSMS_LOGIN")
	if login == "" {
		log.Fatal("Unauthorized: REDSMS_LOGIN env is not set")
	}
	apiKey := os.Getenv("REDSMS_APIKEY")
	if apiKey == "" {
		log.Fatal("Unauthorized: REDSMS_APIKEY env is not set")
	}
	ctx := context.Background()
	tp := redsms.SimpleAuthTransport{
		Login:  login,
		APIKey: apiKey,
	}
	client := redsms.NewClient(tp.Client())
	info, _, err := client.Client.GetInfo(ctx)
	if err != nil {
		log.Fatal(err)
	}

	st := "active"
	if !info.Active {
		st = "inactive"
	}
	log.Printf("Client: %s, status: %s", info.Login, st)
	log.Printf("Balance: %.2f RUB, overdraft: %.2f RUB",
		info.Balance, conv(info.Overdraft))
}

func conv(f *float64) float64 {
	if f == nil {
		return 0.0
	}
	return *f
}
