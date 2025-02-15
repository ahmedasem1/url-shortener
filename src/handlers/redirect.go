package handlers

import (
	"context"
	"time"
	"url-shortener/src/storage"
	"url-shortener/src/workers"

	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/redis/go-redis/v9"
)

type ShortenRequest struct {
	URL string `json:"url" binding:"required"`
}

// @Summary Redirect to original URL
// @Description Redirects a short URL to its original long URL.
// @Tags URL
// @Param shortID path string true "Short URL identifier"
// @Success 302 "Redirects to original URL"
// @Failure 404 {object} map[string]string "Short ID not found"
// @Router /{shortID} [get]
func RedirectURL(c *gin.Context, logChan chan<- workers.AnalyticsLog) {
	shortID := c.Param("shortID")

	url, err := storage.GetURL(c, shortID)
	if err != nil {
		c.JSON(404, gin.H{"error": "URL not found"})
		return
	}

	logChan <- workers.AnalyticsLog{
		ShortID:   shortID,
		Timestamp: time.Now(),
		IP:        c.ClientIP(),
	}

	c.Redirect(302, url)
}

// @Summary Shorten a URL
// @Description Generates a short URL for a given long URL.
// @Tags URL
// @Accept json
// @Produce json
// @Param request body handlers.ShortenRequest true "URL to shorten"
// @Success 200 {object} map[string]string "Shortened URL"
// @Failure 400 {object} map[string]string "Invalid input"
// @Router /shorten [post]
func ShortenURL(c *gin.Context, rdb *redis.Client) {
	var request ShortenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	shortID, err := gonanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 8)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate short ID"})
		return
	}

	ctx := context.Background()
	err = rdb.Set(ctx, shortID, request.URL, 30*24*time.Hour).Err()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to store URL"})
		return
	}

	c.JSON(200, gin.H{"short_url": shortID})
}
