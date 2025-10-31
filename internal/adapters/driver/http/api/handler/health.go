package handler

import (
	"encoding/json"
	"net/http"
)

type health interface {
	CheckHealth(w http.ResponseWriter, r *http.Request)
}

func (h *handler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	health := h.use.CheckHealth(r.Context())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(health); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
