package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"url-shortener/src/handlers"
	"url-shortener/src/storage"
	"url-shortener/src/workers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	setupRedis()
	r := gin.Default()

	r.POST("/shorten", func(c *gin.Context) {
		handlers.ShortenURL(c, storage.RedisClient)
	})

	analyticsChan := make(chan workers.AnalyticsLog, 10)
	r.GET("/:shortID", func(c *gin.Context) {
		handlers.RedirectURL(c, analyticsChan)
	})

	return r
}

func TestShortenURL(t *testing.T) {
	router := setupTestRouter()

	// Test request
	requestBody := map[string]string{"url": "https://example.com"}
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotEmpty(t, response["short_url"])
}

func TestRedirectURL(t *testing.T) {
	router := setupTestRouter()
	shortID := "test123"

	storage.RedisClient.Set(context.Background(), shortID, "https://example.com", 30*24*time.Hour)

	req, _ := http.NewRequest("GET", "/"+shortID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "https://example.com", w.Header().Get("Location"))
}

func TestRateLimitExceeded(t *testing.T) {
	router := setupRateLimitTestRouter()
	ip := "127.0.0.1"

	// Simulate excessive requests
	for i := 0; i < 15; i++ {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = ip + ":12345"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = ip + ":12345"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)

}
