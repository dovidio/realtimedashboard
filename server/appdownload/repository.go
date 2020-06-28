package appdownload

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository provides access to the appdownload storage
type Repository interface {
	// GetAll returns all appdownloads saved in storage
	GetAll() []AppDownload

	// Add saves a given appdownload to the repository
	Add(AppDownload) error
}

type MongoRepository struct {
	client *mongo.Client
}

func NewMongoRepository(client *mongo.Client) Repository {
	return &MongoRepository{client: client}
}

func (m *MongoRepository) GetAll() []AppDownload {
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

func (m *MongoRepository) Add(appDownload AppDownload) error {
	collection := m.client.Database("appdownloads").Collection("appdownloads")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, appDownload)

	if err != nil {
		return err
	}

	return nil
}
