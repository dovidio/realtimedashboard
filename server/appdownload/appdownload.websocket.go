package appdownload

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/websocket"

	"go.mongodb.org/mongo-driver/mongo/readpref"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type message struct {
	Data string `json:"data"`
	Type string `json:"type"`
}

func appDownloadSocket(ws *websocket.Conn) {

	go func(c *websocket.Conn) {
		for {
			var msg message
			if err := websocket.JSON.Receive(ws, &msg); err != nil {
				log.Println(err)
				break
			}

			fmt.Printf("recieved message %s\n", msg.Data)
		}

	}(ws)

	for {
		var appdownloads AppDownload
		if err := websocket.JSON.Send(ws, appdownloads); err != nil {
			break
		}
	}
}

func watchAppDownloads() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017/").SetDirect(true))
	if err != nil {
		log.Fatalf("Failed building mongo client: %v", err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Nearest())
	if err != nil {
		log.Fatalf("failed connecting to mongo: %v", err)
	}

	collection := client.Database("appdownloads").Collection("appdownloads")
	var pipeline mongo.Pipeline // set up pipeline
	streamCur, err := collection.Watch(context.Background(), pipeline, options.ChangeStream())
	if err != nil {
		log.Fatalf("Error getting streaming cursor: %v", err)
	}
	for streamCur.Next(context.Background()) {
		var result bson.M
		streamCur.Decode(&result)

		_, found := result["fullDocument"]
		if !found {
			log.Fatalf("Cannnot find full document: %v", err)
		}

	}
}
