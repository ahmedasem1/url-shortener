package workers

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type AnalyticsLog struct {
	ShortID   string
	Timestamp time.Time
	IP        string
}

func AnalyticsWorker(client *mongo.Client, logChan <-chan AnalyticsLog) {
	collection := client.Database("analytics").Collection("logs")
	for log := range logChan {
		_, err := collection.InsertOne(context.Background(), log)
		if err != nil {
			// Handle error (e.g., log to a file or retry)
		}
	}
}
