package tests

import (
	"context"
	"testing"
	"url-shortener/src/storage"

	"github.com/stretchr/testify/assert"
)

func TestRedisStorage(t *testing.T) {
	setupRedis()
	shortID := "test123"
	url := "https://example.com"

	// Save URL
	err := storage.SaveURL(context.Background(), shortID, url)
	assert.NoError(t, err)

	// Retrieve URL
	retrievedURL, err := storage.GetURL(context.Background(), shortID)
	assert.NoError(t, err)
	assert.Equal(t, url, retrievedURL)
}
