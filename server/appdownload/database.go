package appdownload

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// GetMongoClient does the initial connection to mongo db and returns the client
func GetMongoClient() *mongo.Client {
	ctxWithDeadline, _ := context.WithTimeout(context.Background(), 2*time.Second)
	client, err := mongo.Connect(ctxWithDeadline, options.Client().ApplyURI("mongodb://localhost:27017/").SetDirect(true))
	if err != nil {
		log.Fatalf("Could not get client: %v", err)
	}
	if err = client.Ping(ctxWithDeadline, readpref.Nearest()); err != nil {
		log.Fatalf("failed connecting to mongo: %v", err)
	}

	return client
}
