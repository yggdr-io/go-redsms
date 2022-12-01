package main

import (
	"context"
	"log"
	"os"

	"github.com/gurza/go-redsms/redsms"
)

func main() {
	tp := redsms.SimpleAuthTransport{
		Login:  os.Getenv("LOGIN"),
		APIKey: os.Getenv("APIKEY"),
	}
	c := redsms.NewClient(tp.Client())
	info, _, err := c.Client.GetInfo(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Client: %s", info.Login)
	log.Printf("Balance: %.2f RUB, overdraft: %.2f RUB",
		info.Balance, conv(info.Overdraft))
}

func conv(f *float64) float64 {
	if f == nil {
		return 0.0
	}
	return *f
}
