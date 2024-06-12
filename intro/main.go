package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// read number from redis
	num, err := client.Get(ctx, "number").Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		panic(err)
	}

	log.Printf("Number: %d\n", num)

	// increment number
	err = client.Incr(ctx, "number").Err()
	if err != nil {
		panic(err)
	}

	log.Println("Number incremented")

	// set string to redis with expiration time
	err = client.Set(ctx, "message", "Hello, World!", 2*time.Second).Err()
	if err != nil {
		panic(err)
	}

	log.Println("Sleeping for 3 seconds")
	time.Sleep(3 * time.Second)

	// read string from redis
	message, err := client.Get(ctx, "message").Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		panic(err)
	}

	log.Printf("Message: %s\n", message)
}
