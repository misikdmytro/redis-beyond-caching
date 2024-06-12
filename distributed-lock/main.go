package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/misikdmytro/redis-beyond-caching/distributed-lock/lock"
	"github.com/redis/go-redis/v9"
)

func main() {
	log.Println("starting application...")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	log.Println("acquiring lock...")
	l := lock.NewLock(client, "mylock", 100*time.Millisecond, 5*time.Second)

	err := l.Lock(ctx)
	if err != nil {
		log.Fatalf("could not acquire lock: %v", err)
	}
	defer l.Unlock(ctx)

	log.Println("lock acquired. sleep for 15 seconds...")
	select {
	case <-time.After(15 * time.Second):
		log.Println("wake up")
	case <-ctx.Done():
		log.Println("context cancelled")
	}

	log.Println("done")
}
