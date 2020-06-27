package appdownload

import (
	"log"

	"golang.org/x/net/websocket"
)

// DownloadHandler gets called every time a new download has been observed
type DownloadHandler interface {
	OnNewDownload(AppDownload)
}

type websocketDownloadHandler struct {
	ws              *websocket.Conn
	databaseWatcher DatabaseWatchHandler
}

type message struct {
	Data string `json:"data"`
	Type string `json:"type"`
}

func (w *websocketDownloadHandler) webappDownloadSocket(ws *websocket.Conn) {

	w.ws = ws
	myID := w.databaseWatcher.RegisterHandler(w)
	log.Printf("My id is %d", myID)
	for {
		var msg message
		if err := websocket.JSON.Receive(ws, &msg); err != nil {
			log.Printf("Error while reading message: %v", err)
			break
		}

		log.Printf("received message %s\n", msg.Data)
	}
	w.databaseWatcher.UnregisterHandler(myID)
}

func (w *websocketDownloadHandler) OnNewDownload(a AppDownload) {
	if err := websocket.JSON.Send(w.ws, a); err != nil {
		log.Printf("Error while trying to send update to websocket: %v", err)
	}
}
