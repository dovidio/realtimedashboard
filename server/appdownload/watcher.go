package appdownload

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DatabaseWatchHandler watches for changes in a database and handles updates calling DownloadWatcher.OnNewDownload
type DatabaseWatchHandler interface {
	RegisterObserver(AppDownloadObserver) uuid.UUID
	UnregisterObserver(u uuid.UUID)
	WatchAppDownloads()
}

func NewMongoWatchHandler(c *mongo.Client) DatabaseWatchHandler {
	return &MongoDbWatchHandler{observers: make(map[uuid.UUID]AppDownloadObserver, 0), client: c}
}

type MongoDbWatchHandler struct {
	observers map[uuid.UUID]AppDownloadObserver
	mut       sync.Mutex
	client    *mongo.Client
}

func (m *MongoDbWatchHandler) RegisterObserver(observer AppDownloadObserver) uuid.UUID {
	m.mut.Lock()
	defer m.mut.Unlock()
	uuid := uuid.New()
	m.observers[uuid] = observer

	return uuid
}

func (m *MongoDbWatchHandler) UnregisterObserver(u uuid.UUID) {
	delete(m.observers, u)
}

func (m *MongoDbWatchHandler) WatchAppDownloads() {
	collection := m.client.Database("appdownloads").Collection("appdownloads")
	var pipeline = mongo.Pipeline{
		{{"$project", bson.D{{"operationType", 0}, {"ns", 0}, {"documentKey", 0}, {"clusterTime", 0}}}},
	}
	ctx := context.Background()
	streamCur, err := collection.Watch(ctx, pipeline, options.ChangeStream().SetFullDocument(options.UpdateLookup))
	if err != nil {
		log.Fatalf("Error getting streaming cursor: %v", err)
	}

	for streamCur.Next(ctx) {
		var result bson.M
		streamCur.Decode(&result)

		fullDocument, found := result["fullDocument"]
		if !found {
			log.Fatalf("Cannnot find full document: %v", err)
		}
		appDownload, err := extractAppDownload(fullDocument)
		if err != nil {
			fmt.Print(err)
		}

		m.mut.Lock()
		for _, observer := range m.observers {
			observer.OnNewAppDownload(appDownload)
		}
		m.mut.Unlock()
	}
}

func extractAppDownload(m interface{}) (AppDownload, error) {

	var appDownload AppDownload
	b, ok := m.(bson.M)
	if !ok {
		return appDownload, errors.New("Could not deserialize document")
	}

	mapstructure.Decode(b, &appDownload)

	return appDownload, nil
}

// AppDownloadObserver gets called every time a new download has been observed
type AppDownloadObserver interface {
	OnNewAppDownload(AppDownload)
}
