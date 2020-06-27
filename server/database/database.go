package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client exposes a mongo client
var Client *mongo.Client

// SetupDatabase does the initial connection to mongo db and initialize the client
func SetupDatabase() {
	var err error
	Client, err = mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017/").SetDirect(true))

	if err != nil {
		log.Fatalf("Could not get client: %v", err)
	}
}
