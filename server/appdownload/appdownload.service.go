package appdownload

import (
	"encoding/json"
	"log"
	"net/http"
	"realtimedashboard/appdownload/cors"
)

const appDownloadsBasePath = "/appdownloads"

// SetupRoutes registers the routes for app downloads
func SetupRoutes() {
	handler := &Handler{}
	http.Handle(appDownloadsBasePath, cors.Middleware(handler))
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
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
