package appdownload

import (
	"encoding/json"
	"log"
	"net/http"
)

// Handler handles request for appdownloads
type Handler struct {
	repository Repository
}

func NewHandler(repository Repository) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// GetAll returns appdownloads as json
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(h.repository.GetAll())
	if err != nil {
		log.Panicf("Could not marshall to json: %v", err)
	}
	w.Write(b)
}
