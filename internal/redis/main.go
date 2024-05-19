package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func main() {
	// Create a new Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Replace with your Redis server address
		Password: "",               // Replace with your Redis server password, if any
		DB:       0,                // Replace with your Redis database number
	})

	ctx := context.Background()
	// Ping the Redis server to check the connection
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis:", err)
		return
	}

	fmt.Println("Connected to Redis:", pong)

	// Close the Redis client when you're done
	err = client.Close()
	if err != nil {
		fmt.Println("Failed to close Redis client:", err)
	}
}
