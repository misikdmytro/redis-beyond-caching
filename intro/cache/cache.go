package cache

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/redis/go-redis/v9"
)

type CacheValue struct {
	Number  int
	Message string
}

func CacheAside(ctx context.Context) (CacheValue, error) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// 1. read cache from redis
	result, err := client.Get(ctx, "value").Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return CacheValue{}, err
	}

	if result != "" {
		// cache found
		var cv CacheValue
		if err := json.Unmarshal([]byte(result), &cv); err != nil {
			return CacheValue{}, err
		}

		return cv, nil
	}

	// 2. if cache not found, read from data source
	cv, err := readFromDataSource(ctx)
	if err != nil {
		return CacheValue{}, err
	}

	// 3. set cache to redis
	data, err := json.Marshal(cv)
	if err != nil {
		return CacheValue{}, err
	}

	if err := client.Set(ctx, "value", data, 0).Err(); err != nil {
		return CacheValue{}, err
	}

	return cv, nil
}

func readFromDataSource(ctx context.Context) (CacheValue, error) {
	// read from data source
	return CacheValue{
		Number:  42,
		Message: "Hello, World!",
	}, nil
}
