package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimit applies rate limiting using Redis
func RateLimit(rdb *redis.Client, limit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := "rate_limit:" + ip

		// Increment request count
		count, err := rdb.Incr(c, key).Result()
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "Rate limit error"})
			return
		}

		// Set expiration on first request
		if count == 1 {
			rdb.Expire(c, key, time.Minute)
		}

		// Check if rate limit exceeded
		if count > int64(limit) {
			c.AbortWithStatusJSON(429, gin.H{"error": "Rate limit exceeded"})
			return
		}

		c.Next()
	}
}
