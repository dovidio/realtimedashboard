package appdownload

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// SetupRoutes registers the routes for app downloads
func SetupRoutes() {
	handler := &Handler{}
	http.Handle("/appdownloads", handler)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Handler serves appdownload metadata
type Handler struct{}

func (a *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("called")
	b, err := json.Marshal(getAppDownloadList())
	if err != nil {
		log.Panicf("Could not marshall to json: %v", err)
	}

	w.Write(b)
}
