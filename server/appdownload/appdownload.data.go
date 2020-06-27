package appdownload

import (
	"context"
	"log"
	"realtimedashboard/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var MockData = []interface{}{
	AppDownload{
		Longitude:    12.377281,
		Latitude:     42.000325,
		AppID:        "IOS_ALERT",
		DownloadedAt: 1593181859189,
	},
	AppDownload{
		Longitude:    11.377481,
		Latitude:     41.000325,
		AppID:        "IOS_MATE",
		DownloadedAt: 1593183859189,
	},
	AppDownload{
		Longitude:    22.377281,
		Latitude:     44.000325,
		AppID:        "ANDROID_ALERT",
		DownloadedAt: 1593184859189,
	},
	AppDownload{
		Longitude:    12.377281,
		Latitude:     42.000325,
		AppID:        "ANDROID_MATE",
		DownloadedAt: 1593189853189,
	},
}

func getAppDownloadList() []AppDownload {
	collection := database.Client.Database("appdownloads").Collection("appdownloads")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatalf("There was an error while retrieving app downloads: %v", err)
	}

	appDownloads := make([]AppDownload, 0)
	for cur.Next(ctx) {

		var appDownload AppDownload
		cur.Decode(&appDownload)
		if err != nil {
			log.Printf("Error while marshalling appdownload: %v", err)
		}

		appDownloads = append(appDownloads, appDownload)
	}

	// collection

	return appDownloads
}
