package main

import (
	"log"
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
	// Initialize Redis
	storage.InitRedis("redis:6379")

	// Initialize MongoDB
	mongoClient, err := storage.InitMongoDB("mongodb://mongo:27017")
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Create a channel for analytics logging
	analyticsChan := make(chan workers.AnalyticsLog)

	// Start worker goroutines
	go workers.AnalyticsWorker(mongoClient, analyticsChan)

	// Set up Gin router
	r := gin.Default()
	r.Use(middleware.RateLimit(storage.RedisClient))

	r.POST("/shorten", func(c *gin.Context) {
		handlers.ShortenURL(c, storage.RedisClient)
	})

	r.GET("/:shortID", func(c *gin.Context) {
		handlers.RedirectURL(c, analyticsChan)
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Server started on port 8080")
	r.Run(":8080")
}
