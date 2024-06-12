package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/misikdmytro/redis-beyond-caching/rate-limiter/middleware"
	"github.com/redis/go-redis/v9"
)

func main() {
	app := fiber.New()

	app.Use(middleware.NewRateLimiter(middleware.RateLimiterOptions{
		Redis: redis.Options{
			Addr: "localhost:6379",
		},
		Limit:     5,
		KeyPrefix: "rate-limiter",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":3000")
}
