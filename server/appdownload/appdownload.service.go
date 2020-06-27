package appdownload

import (
	"encoding/json"
	"log"
	"net/http"
	"realtimedashboard/appdownload/cors"
	"time"

	"golang.org/x/net/websocket"
)

const appDownloadsBasePath = "/appdownloads"
const appDownloadsWebsocketBasePath = "/appdownloadssocket"

// SetupRoutes registers the routes for app downloads
func SetupRoutes() {
	handler := &Handler{}
	http.Handle(appDownloadsBasePath, cors.Middleware(handler))
	dbWatcher := &mongoDbWatchHandler{watchHandlers: make([]DownloadHandler, 0)}
	go dbWatcher.WatchAppDownloads()
	wsHandler := &websocketDownloadHandler{databaseWatcher: dbWatcher}
	http.Handle(appDownloadsWebsocketBasePath, websocket.Handler(wsHandler.webappDownloadSocket))

	go generateDataPeriodically()

	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func generateDataPeriodically() {
	for {
		time.Sleep(1 * time.Second)
		insertRandomDownload()
	}
}

// Handler serves appdownload metadata
type Handler struct{}

func (a *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		b, err := json.Marshal(getAppDownloadList())
		if err != nil {
			log.Panicf("Could not marshall to json: %v", err)
		}

		w.Write(b)
		return
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
