package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type RateLimiterOptions struct {
	Redis     redis.Options
	Limit     int
	KeyPrefix string
}

func NewRateLimiter(options RateLimiterOptions) fiber.Handler {
	client := redis.NewClient(&options.Redis)

	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		now := time.Now()

		key := fmt.Sprintf("%s:%s:%s", options.KeyPrefix, c.IP(), now.Format("15:04"))

		// Increment the key
		value, err := client.Incr(ctx, key).Result()
		if err != nil {
			return err
		}

		// If the value is greater than the limit, return an error
		if value > int64(options.Limit) {
			return fiber.ErrTooManyRequests
		}

		return c.Next()
	}
}
