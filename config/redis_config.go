package configs

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

var (
	ctx    = context.Background()
	client *redis.Client
)

func init() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL environment variable is not set.")
	}

	options, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatal("Error parsing Redis URL:", err)
	}

	client = redis.NewClient(options)

	// Check if the Redis connection is successful
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}
	fmt.Println("Connected to Redis:", pong)
}

func GetString(key string) (string, error) {
	result, err := client.Get(ctx, key).Result()
	if err != nil {
		return "null", err
	}
	return result, nil
}

func SetString(key string, value string) error {
	err := client.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}
