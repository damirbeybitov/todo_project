package redis

import (
	"context"
	"fmt"

	"github.com/damirbeybitov/todo_project/internal/log"
	"github.com/redis/go-redis/v9"
)

// NewClient initializes a new Redis client
func NewClient(addr, password string, db int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.ErrorLogger.Fatalf("failed to connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis:", pong)

	return client
}
