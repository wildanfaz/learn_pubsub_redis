package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func main() {
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// create a channel for receiving messages
	messages := make(chan *redis.Message)

	// create a Goroutine to handle message subscriptions
	go func() {
		pubsub := client.Subscribe(ctx, "mychannel")
		defer pubsub.Close()

		// listen for messages on the channel
		ch := pubsub.Channel()
		for msg := range ch {
			messages <- msg
		}
	}()

	// create a Goroutine to handle message publishing
	go func() {
		for {
			var x string
			fmt.Scan(&x)
			x = "flovvint:" + x
			err := client.Publish(ctx, "mychannel", x).Err()
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	// receive messages from the channel
	for msg := range messages {
		fmt.Println(msg.Payload)
	}
}
