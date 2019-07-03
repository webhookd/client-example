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
	for {

		msg := webhookMessage{Name: "test message", Time: time.Now(), Uuid: uuid.NewV4().String()}
		jsnBody, err := json.Marshal(msg)

		webhooks := webhookd.NewWebhookdApiClient(webhookdConnection)
		provider := &webhookd.Provider{Id: "chargehive"}

		subscription := &webhookd.Subscription{}
		subscription.Key = "testing"

		message := &webhookd.Message{
			Id:     "example-id",
			Method: webhookd.HTTP_METHOD_POST,
			Body:   string(jsnBody),
			Query: map[string]string{
				"token": "accessToken",
			},
			Headers: map[string]string{
				"x-example-url": "http://originated.at.localhost",
				"content-type":  "application/json",
			},
		}

		if subscription.Id == "" && subscription.Key == "" {
			invResp, err := webhooks.SubscriptionCreate(context.Background(), &webhookd.SubscriptionCreateRequest{
				Provider: provider,
				Key:      "client-test" + time.Now().Format("Jan-2-15-04-05"),
				Name:     "Client Test " + time.Now().Format(time.RFC822),
			})

			if err != nil {
				log.Fatal(err)
			}

			subscription = invResp.Subscription

			accept, err := webhooks.SubscriptionAccept(context.Background(), &webhookd.SubscriptionAcceptRequest{
				Subscription: subscription,
				Consumer:     &webhookd.Consumer{Id: "random-consumer-id"},
				Token:        invResp.Token,
			})

			if err != nil || accept.Received {
				log.Fatal(err, accept)
			}
		}

		resp, err := webhooks.Transmit(context.Background(), &webhookd.TransmitRequest{
			Provider:     provider,
			Subscription: subscription,
			Message:      message,
		})
		log.Print(resp, err)
		time.Sleep(time.Second * 20)
	}
}
