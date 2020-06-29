package main

import (
	"log"
	"net/http"
	"os"
	"realtimedashboard/appdownload"
	"realtimedashboard/cors"
	"realtimedashboard/db"
	"strconv"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

func main() {
	db := db.NewClient().Database("appdownloads")
	repository := appdownload.NewMongoRepository(db)
	dbWatcher := appdownload.NewMongoWatchHandler(db, make(map[uuid.UUID]appdownload.Observer, 0))

	handler := appdownload.NewHandler(repository)
	wsHandler := appdownload.NewWebsocketHandler(dbWatcher)

	http.Handle("/appdownloads", cors.Middleware(http.HandlerFunc(handler.Handle)))
	http.Handle("/appdownloadssocket", websocket.Handler(wsHandler.Websocket))

	go dbWatcher.WatchAppDownloads()

	if period := shouldGenerateData(); period > 0 {
		quit := make(chan struct{})
		go appdownload.GenerateData(time.Duration(period)*time.Millisecond, repository, quit)
	}

	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func shouldGenerateData() int64 {
	generateData, err := strconv.ParseBool(os.Getenv("GENERATE_DATA"))
	if err == nil && generateData {
		period, err := strconv.ParseInt(os.Getenv("GENERATE_DATA_PERIOD"), 10, 32)
		if err == nil {
			return period
		}
	}

	return -1
}
