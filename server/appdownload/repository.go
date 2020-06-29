package appdownload

import (
	"context"
	"log"
	"realtimedashboard/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// Repository provides access to the appdownload storage
type Repository interface {
	// GetAll returns all appdownloads saved in storage
	GetAll() []AppDownload

	// Add saves a given appdownload to the repository
	Add(AppDownload) error
}

// MongoRepository provides access to the appdownload mongo storage
type MongoRepository struct {
	db db.DatabaseHelper
}

// NewMongoRepository creates a new MongoRepository
func NewMongoRepository(db db.DatabaseHelper) Repository {
	return &MongoRepository{db: db}
}

// GetAll returns all appdownloads stored in mongo
func (m *MongoRepository) GetAll() []AppDownload {
	collection := m.db.Collection("appdownloads")

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

// Add adds a new appdownload to mongo
func (m *MongoRepository) Add(appDownload AppDownload) error {
	collection := m.db.Collection("appdownloads")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, appDownload)

	if err != nil {
		return err
	}

	return nil
}
