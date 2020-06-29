package main

import (
	"log"
	"net/http"
	"realtimedashboard/appdownload"
	"realtimedashboard/cors"
	"time"

	"golang.org/x/net/websocket"
)

func main() {
	client := appdownload.GetMongoClient()
	repository := appdownload.NewMongoRepository(client)
	dbWatcher := appdownload.NewMongoWatchHandler(client)

	handler := appdownload.NewHandler(repository)
	wsHandler := appdownload.NewWebsocketHandler(dbWatcher)

	http.Handle("/appdownloads", cors.Middleware(http.HandlerFunc(handler.Handle)))
	http.Handle("/appdownloadssocket", websocket.Handler(wsHandler.Websocket))

	go dbWatcher.WatchAppDownloads()
	go appdownload.GenerateData(1*time.Second, repository)

	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
