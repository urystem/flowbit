package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

type modeResponse struct {
	Mode      string    `json:"mode"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type mode interface {
	SwitchToTest(w http.ResponseWriter, r *http.Request)
	SwitchToLive(w http.ResponseWriter, r *http.Request)
}

func (h *handler) SwitchToTest(w http.ResponseWriter, r *http.Request) {
	h.use.SwitchToTest()

	resp := modeResponse{
		Mode:      "test",
		Message:   "System switched to test mode. Fetching data from generator.",
		Timestamp: time.Now().UTC(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *handler) SwitchToLive(w http.ResponseWriter, r *http.Request) {
	h.use.SwitchToLive()

	resp := modeResponse{
		Mode:      "live",
		Message:   "System switched to live mode. Fetching data from exchanges.",
		Timestamp: time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
