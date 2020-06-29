package appdownload

import (
	"log"

	"golang.org/x/net/websocket"
)

// WebsocketHandler handles the websocket connection and register itself to the DatabaseWatchHandler
type WebsocketHandler struct {
	ws              *websocket.Conn
	databaseWatcher DatabaseWatchHandler
}

type message struct {
	Data string `json:"data"`
	Type string `json:"type"`
}

// NewWebsocketHandler requires a databaseWatchHandler for registering to changes
func NewWebsocketHandler(databaseWatchHandler DatabaseWatchHandler) *WebsocketHandler {
	return &WebsocketHandler{databaseWatcher: databaseWatchHandler}
}

// Websocket streams appdownloads to the websocket connection
func (w *WebsocketHandler) Websocket(ws *websocket.Conn) {
	w.ws = ws
	myID := w.databaseWatcher.RegisterObserver(w)
	log.Printf("My id is %d", myID)
	for {
		var msg message
		if err := websocket.JSON.Receive(ws, &msg); err != nil {
			log.Printf("Error while reading message: %v", err)
			break
		}

		log.Printf("received message %s\n", msg.Data)
	}
	w.databaseWatcher.UnregisterObserver(myID)
}

// OnNewAppDownload send the new app downoads to the websocket
func (w *WebsocketHandler) OnNewAppDownload(a AppDownload) {
	if err := websocket.JSON.Send(w.ws, a); err != nil {
		log.Printf("Error while trying to send update to websocket: %v", err)
	}
}
