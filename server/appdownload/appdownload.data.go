package appdownload

import (
	"context"
	"log"
	"math/rand"
	"os"
	"realtimedashboard/database"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Longitude    float64            `bson:"longitude" json:"longitude"`
	Latitude     float64            `bson:"latitude" json:"latitude"`
	AppID        string             `bson:"app_id" json:"app_id"`
	DownloadedAt int64              `bson:"downloaded_at" json:"downloaded_at"`
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
	appDownload.Latitude = rand.Float64()*180.0 - 90.0
	appDownload.Longitude = rand.Float64()*180.0 - 90.0
	appDownload.DownloadedAt = time.Now().Unix()

	_, err := collection.InsertOne(ctx, appDownload)

	if err != nil {
		log.Panicf("Error while trying to insert an appdownload: %v", err)
	}
}

// DatabaseWatchHandler watches for changes in a database and handles updates calling DownloadWatcher.OnNewDownload
type DatabaseWatchHandler interface {
	RegisterHandler(d DownloadHandler) int
	UnregisterHandler(id int)
	WatchAppDownloads()
}

type mongoDbWatchHandler struct {
	watchHandlers []DownloadHandler
	mut           sync.Mutex
}

func (m *mongoDbWatchHandler) RegisterHandler(d DownloadHandler) int {
	m.mut.Lock()
	defer m.mut.Unlock()
	m.watchHandlers = append(m.watchHandlers, d)
	return len(m.watchHandlers) - 1
}

func (m *mongoDbWatchHandler) UnregisterHandler(id int) {
	m.watchHandlers = append(m.watchHandlers[:id], m.watchHandlers[id+1:]...)
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
	var pipeline mongo.Pipeline
	ctx := context.Background()
	streamCur, err := collection.Watch(ctx, pipeline, options.ChangeStream())
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
		appDownload := extractAppDownload(fullDocument)

		m.mut.Lock()
		for _, watcher := range m.watchHandlers {
			watcher.OnNewDownload(appDownload)
		}
		m.mut.Unlock()
	}
}

func extractAppDownload(m interface{}) AppDownload {
	var appDownload AppDownload
	appDownload, ok := m.(AppDownload)
	if !ok {
		log.Printf("got data of type %T but wanted int", m)
		os.Exit(1)
	}

	return appDownload
}
