package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"url-shortener/src/middleware"
	"url-shortener/src/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Initialize Redis once for all tests
func setupRedis() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}
	storage.InitRedis(redisAddr)

	// Flush Redis before running tests
	storage.RedisClient.FlushDB(context.Background())
}

// Setup router for rate limiting tests
func setupRateLimitTestRouter() *gin.Engine {
	setupRedis()
	r := gin.Default()
	r.Use(middleware.RateLimit(storage.RedisClient))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	return r
}

func TestRateLimit(t *testing.T) {
	router := setupRateLimitTestRouter()
	ip := "127.0.0.1"

	// Simulate 10 requests (allowed)
	for i := 0; i < 10; i++ {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = ip + ":12345"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// 11th request (should be blocked)
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = ip + ":12345"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}

func TestRateLimitReset(t *testing.T) {
	router := setupRateLimitTestRouter()
	ip := "127.0.0.1"
	redisKey := "rate_limit:" + ip

	// Send 10 requests (allowed)
	for i := 0; i < 10; i++ {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = ip + ":12345"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// Manually reset rate limit key in Redis
	storage.RedisClient.Del(context.Background(), redisKey)

	// Wait for rate limit expiry if necessary
	time.Sleep(1 * time.Second)

	// Request after reset (should be allowed)
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = ip + ":12345"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
