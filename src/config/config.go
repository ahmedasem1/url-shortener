package config

import (
	"log"
	"os"
	"strconv"
)

// Config holds application configuration values
type Config struct {
	RedisAddr  string
	MongoURI   string
	ServerPort string
	RateLimit  int
}

// LoadConfig reads environment variables or sets defaults
func LoadConfig() *Config {
	rateLimit, err := strconv.Atoi(getEnv("RATE_LIMIT", "10"))
	if err != nil {
		log.Fatalf("Invalid RATE_LIMIT value: %v", err)
	}

	return &Config{
		RedisAddr:  getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		MongoURI:   getEnv("MONGO_URI", "mongodb://localhost:27017"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		RateLimit:  rateLimit,
	}
}

// getEnv fetches an environment variable or returns a default value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
