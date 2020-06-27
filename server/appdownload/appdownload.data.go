package appdownload

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"realtimedashboard/database"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
)

// AppID identifies the os and type of the application
type AppID string

// AppNames hold a list of possible app names
var AppNames = [...]string{
	"IOS_ALERT",
	"IOS_MATE",
	"IOS_E4",
	"ANDROID_ALERT",
	"ANDROID_MATE",
	"ANDOID_E4",
}

// AppDownload contains metadata about downloads
type AppDownload struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" mapstructure:"_id"`
	Longitude    float64            `bson:"longitude" json:"longitude"`
	Latitude     float64            `bson:"latitude" json:"latitude"`
	AppID        string             `bson:"app_id" json:"app_id" mapstructure:"app_id"`
	DownloadedAt int64              `bson:"downloaded_at" json:"downloaded_at" mapstructure:"downloaded_at"`
}

func getAppDownloadList() []AppDownload {
	collection := database.Client.Database("appdownloads").Collection("appdownloads")

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

func insertRandomDownload() {
	collection := database.Client.Database("appdownloads").Collection("appdownloads")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var appDownload AppDownload
	appDownload.AppID = AppNames[rand.Int31n(int32(len(AppNames)))]
	appDownload.Latitude = rand.Float64()*10 + 40.0
	appDownload.Longitude = rand.Float64()*10 - 10.0
	appDownload.DownloadedAt = time.Now().Unix()

	_, err := collection.InsertOne(ctx, appDownload)

	if err != nil {
		log.Panicf("Error while trying to insert an appdownload: %v", err)
	}
}

// DatabaseWatchHandler watches for changes in a database and handles updates calling DownloadWatcher.OnNewDownload
type DatabaseWatchHandler interface {
	RegisterHandler(d DownloadHandler) uuid.UUID
	UnregisterHandler(u uuid.UUID)
	WatchAppDownloads()
}

type mongoDbWatchHandler struct {
	watchHandlers map[uuid.UUID]DownloadHandler
	mut           sync.Mutex
}

func (m *mongoDbWatchHandler) RegisterHandler(d DownloadHandler) uuid.UUID {
	m.mut.Lock()
	defer m.mut.Unlock()
	uuid, _ := uuid.New()
	m.watchHandlers[uuid] = d

	return uuid
}

func (m *mongoDbWatchHandler) UnregisterHandler(u uuid.UUID) {
	delete(m.watchHandlers, u)
}

func (m *mongoDbWatchHandler) WatchAppDownloads() {
	log.Print("initialized watchAppDownloads")
	ctxWithDeadline, _ := context.WithTimeout(context.Background(), 2*time.Second)
	client, err := mongo.Connect(ctxWithDeadline, options.Client().ApplyURI("mongodb://localhost:27017/").SetDirect(true))
	if err != nil {
		log.Fatalf("Failed building mongo client: %v", err)
	}

	if err = client.Ping(ctxWithDeadline, readpref.Nearest()); err != nil {
		log.Fatalf("failed connecting to mongo: %v", err)
	}

	collection := client.Database("appdownloads").Collection("appdownloads")
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
		for _, watcher := range m.watchHandlers {
			watcher.OnNewDownload(appDownload)
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
