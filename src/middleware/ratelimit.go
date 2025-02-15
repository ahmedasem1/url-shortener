package middleware

import (
	"log"
	"net/http"
	"time"

	"context"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimit(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if redisClient == nil {
			log.Println("Rate limiter disabled: Redis is unavailable")
			c.Next()
			return
		}

		ctx := context.Background()
		ip := c.ClientIP()
		key := "rate_limit:" + ip

		count, err := redisClient.Incr(ctx, key).Result()
		if err != nil {
			log.Println("Redis error in rate limiter:", err)
			c.Next()
			return
		}

		if count == 1 {
			redisClient.Expire(ctx, key, 60*time.Second)
		}

		if count > 10 {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()
	}
}
