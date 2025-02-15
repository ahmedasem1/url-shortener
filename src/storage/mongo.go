package storage

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongoDB(uri string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("MongoDB connection failed: %w", err)
	}

	// Verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("MongoDB ping failed: %w", err)
	}

	fmt.Println("Connected to MongoDB")
	return client, nil
}
