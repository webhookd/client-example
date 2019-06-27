package main

import (
	"context"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"github.com/webhookd/api-proto-go/webhookd"
	"google.golang.org/grpc"
	"log"
	"time"
)

type webhookMessage struct {
	Name string
	Time time.Time
	Uuid string
}

func main() {
	webhookdConnection, err := grpc.Dial("localhost:53200", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	msg := webhookMessage{Name: "test message", Time: time.Now(), Uuid: uuid.NewV4().String()}
	jsnBody, err := json.Marshal(msg)

	webhooks := webhookd.NewWebhookdApiClient(webhookdConnection)
	resp, err := webhooks.Transmit(context.Background(), &webhookd.TransmitRequest{
		Provider:     &webhookd.Provider{Id: "provider-id"},
		Subscription: &webhookd.Subscription{Key: "provider-set-subscription-key"},
		Message: &webhookd.Message{
			Method: webhookd.HttpMethod_POST,
			Body:   string(jsnBody),
			Query: map[string]string{
				"token": "accessToken",
			},
			Headers: map[string]string{
				"x-example-url": "http://originated.at.localhost",
				"content-type":  "application/json",
			},
		},
	})
	log.Print(resp, err)
}
