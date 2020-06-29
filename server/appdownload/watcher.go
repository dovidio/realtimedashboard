package appdownload

import (
	"context"
	"errors"
	"fmt"
	"realtimedashboard/db"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DatabaseWatchHandler watches for changes in a database and handles notifies its observers
type DatabaseWatchHandler interface {
	RegisterObserver(Observer) uuid.UUID
	UnregisterObserver(u uuid.UUID)
	WatchAppDownloads()
}

// NewMongoWatchHandler creates a new watch handler
func NewMongoWatchHandler(db db.DatabaseHelper) DatabaseWatchHandler {
	return &MongoDbWatchHandler{observers: make(map[uuid.UUID]Observer, 0), db: db}
}

// MongoDbWatchHandler holds a map of observers, a mutex for synchronizing across different threads,
// and a db helper to subscribe to the db changes
type MongoDbWatchHandler struct {
	observers map[uuid.UUID]Observer
	mut       sync.Mutex
	db        db.DatabaseHelper
}

// RegisterObserver add an observer to the observers map
func (m *MongoDbWatchHandler) RegisterObserver(observer Observer) uuid.UUID {
	m.mut.Lock()
	defer m.mut.Unlock()
	uuid := uuid.New()
	m.observers[uuid] = observer

	return uuid
}

// UnregisterObserver remove the observer from the map
func (m *MongoDbWatchHandler) UnregisterObserver(u uuid.UUID) {
	m.mut.Lock()
	defer m.mut.Unlock()
	delete(m.observers, u)
}

// WatchAppDownloads starts watching changes in db and notifies observers accordingly
func (m *MongoDbWatchHandler) WatchAppDownloads() {
	for {
		if err := m.watchAppDownloads(); err != nil {
			fmt.Println(err)
			time.Sleep(1 * time.Second)
		}
	}
}

func (m *MongoDbWatchHandler) watchAppDownloads() error {
	collection := m.db.Collection("appdownloads")
	var pipeline = mongo.Pipeline{
		{{"$project", bson.D{{"operationType", 0}, {"ns", 0}, {"documentKey", 0}, {"clusterTime", 0}}}},
	}
	ctx := context.Background()
	streamCur, err := collection.Watch(ctx, pipeline, options.ChangeStream().SetFullDocument(options.UpdateLookup))
	if err != nil {
		return fmt.Errorf("Error getting streaming cursor: %v", err)
	}

	for streamCur.Next(ctx) {
		var result bson.M
		streamCur.Decode(&result)

		fullDocument, found := result["fullDocument"]
		if !found {
			return fmt.Errorf("Cannnot find full document: %v", err)
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

	return errors.New("streaming cursor finished")
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

// Observer gets called every time a new download has been observed
type Observer interface {
	OnNewAppDownload(AppDownload)
}
