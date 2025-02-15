package main

import (
	"log"
	"url-shortener/src/config"
	"url-shortener/src/handlers"
	"url-shortener/src/middleware"
	"url-shortener/src/storage"
	"url-shortener/src/workers"

	_ "url-shortener/src/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize Redis
	storage.InitRedis(cfg.RedisAddr)

	// Initialize MongoDB
	mongoClient, err := storage.InitMongoDB(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Create a channel for analytics logging
	analyticsChan := make(chan workers.AnalyticsLog)

	// Start worker goroutines
	go workers.AnalyticsWorker(mongoClient, analyticsChan)

	// Set up Gin router
	r := gin.Default()
	r.Use(middleware.RateLimit(storage.RedisClient, cfg.RateLimit))

	r.POST("/shorten", func(c *gin.Context) {
		handlers.ShortenURL(c, storage.RedisClient)
	})

	r.GET("/:shortID", func(c *gin.Context) {
		handlers.RedirectURL(c, analyticsChan)
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Printf("Server started on port %s", cfg.ServerPort)
	r.Run(":" + cfg.ServerPort)
}
