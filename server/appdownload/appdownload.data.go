package appdownload

import (
	"context"
	"log"
	"math/rand"
	"realtimedashboard/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

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
