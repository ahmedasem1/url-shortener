package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// InitRedis initializes the Redis client
func InitRedis(redisAddr string) {
	if redisAddr == "" {
		redisAddr = os.Getenv("REDIS_ADDR")
		if redisAddr == "" {
			redisAddr = "0.0.0.0:6379"
		}
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	fmt.Println("Connected to Redis:", redisAddr)
}

func GetURL(ctx context.Context, shortID string) (string, error) {
	url, err := RedisClient.Get(ctx, shortID).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("URL not found")
	} else if err != nil {
		return "", err
	}
	return url, nil
}

func SaveURL(ctx context.Context, shortID string, url string) error {
	return RedisClient.Set(ctx, shortID, url, 30*24*time.Hour).Err()
}
