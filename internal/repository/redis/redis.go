package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var Rdb *redis.Client

// InitRedis initializes the Redis connection
func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr: "redis:6379", // Redis address
		DB:   0,            // Redis database to use
	})

	// Check the connection with a Ping command
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Execute the Ping command to test the connection
	if err := Rdb.Ping(ctx).Err(); err != nil {
		panic(fmt.Errorf("could not connect to Redis: %v", err))
	}
}
