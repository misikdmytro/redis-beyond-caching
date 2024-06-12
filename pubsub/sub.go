package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// subscribe to channel
	pubsub := client.Subscribe(ctx, "channel")
	defer pubsub.Close()

	// receive message
	log.Println("waiting for message...")

	for {
		select {
		case <-ctx.Done():
			log.Println("done")
			return
		case msg := <-pubsub.Channel():
			log.Printf("received: %s\n", msg.Payload)
		}
	}
}
