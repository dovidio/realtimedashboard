package main

import (
	"context"
	"fmt"
	"log"
	"realtimedashboard/appdownload"
	"realtimedashboard/database"
	"time"

	"go.mongodb.org/mongo-driver/mongo/readpref"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	database.SetupDatabase()
	appdownload.SetupRoutes()
}

func listenToNewAppDownloads() error {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017/").SetDirect(true))
	if err != nil {
		return fmt.Errorf("Failed building mongo client: %v", err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	err = client.Ping(ctx, readpref.Nearest())
	if err != nil {
		return fmt.Errorf("failed connecting to mongo: %v", err)
	}

	collection := client.Database("appdownloads").Collection("appdownloads")
	go insertMockData(collection)

	return nil
}

func insertMockData(collection *mongo.Collection) {
	for true {
		insertResult, err := collection.InsertMany(context.Background(), appdownload.MockData)

		if err != nil {
			log.Fatal("Error inserting records", err)
		}
		fmt.Println("Inserted element with id", insertResult.InsertedIDs)
		time.Sleep(3 * time.Second)
	}
}
