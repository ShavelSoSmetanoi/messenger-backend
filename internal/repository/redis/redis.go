package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var Rdb *redis.Client

// InitRedis Инициализация Redis
func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr: "redis:6379", // Адрес Redis
		DB:   0,            // Используемая база
	})

	// Проверяем соединение с Redis
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Rdb.Ping(ctx).Err(); err != nil {
		panic(fmt.Errorf("could not connect to Redis: %v", err))
	}
}
