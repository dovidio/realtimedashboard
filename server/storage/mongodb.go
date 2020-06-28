package storage

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoAppDownloadStorage keeps appdownload data in mongo
type MongoAppDownloadStorage struct {
	client *mongo.Client
}

// GetAll returns all appdownloads saved in mongo
func (m *MongoAppDownloadStorage) GetAll() []AppDownload {
	collection := m.client.Database("appdownloads").Collection("appdownloads")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Panicf("There was an error while retrieving app downloads: %v", err)
	}

	appDownloads := make([]AppDownload, 0)
	for cur.Next(ctx) {

		var appDownload AppDownload
		cur.Decode(&appDownload)
		if err != nil {
			log.Panicf("Error while marshalling appdownload: %v", err)
		}

		appDownloads = append(appDownloads, appDownload)
	}

	return appDownloads

}

// Add saves a given appdownload to mongo
func (m *MongoAppDownloadStorage) Add(a AppDownload) error {
	collection := m.client.Database("appdownloads").Collection("appdownloads")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, a)

	if err != nil {
		log.Panicf("Error while trying to insert an appdownload: %v", err)
	}
}
