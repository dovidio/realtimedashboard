package listing

import (
	"encoding/json"
	"log"
	"net/http"
)

func MakeGetAppDownloadsEndpoint(s Service) func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		b, err := json.Marshal(s.GgetAppDownloadList(client))
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
