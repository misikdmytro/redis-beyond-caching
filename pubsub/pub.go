package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// publish message
	msg := uuid.NewString()
	client.Publish(ctx, "channel", msg)

	log.Printf("message published. message: %s\n", msg)
}
