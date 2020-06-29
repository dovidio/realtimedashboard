package appdownload

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// GetMongoClient does the initial connection to mongo db and returns the client
func GetMongoClient() *mongo.Client {
	dbHostname := os.Getenv("DB_HOSTNAME")
	if dbHostname == "" {
		dbHostname = "localhost"
	}
	dbPort, err := strconv.ParseInt(os.Getenv("DB_PORT"), 10, 64)
	if err != nil {
		dbPort = 27017
	}

	dbConnectionString := fmt.Sprintf("mongodb://%s:%d/", dbHostname, dbPort)
	fmt.Println("Trying to connect to ", dbConnectionString)

	ctxWithDeadline, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctxWithDeadline, options.Client().ApplyURI(dbConnectionString).SetDirect(true))
	if err != nil {
		log.Fatalf("Could not get client: %v", err)
	}
	if err = client.Ping(ctxWithDeadline, readpref.Nearest()); err != nil {
		log.Fatalf("failed connecting to mongo: %v", err)
	}

	return client
}
