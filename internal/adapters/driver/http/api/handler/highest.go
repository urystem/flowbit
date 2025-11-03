package handler

import (
	"encoding/json"
	"errors"
	"marketflow/internal/domain"
	"net/http"
	"time"
)

type highest interface {
	GetHighestPriceBySym(w http.ResponseWriter, r *http.Request)
}

func (h *handler) GetHighestPriceBySym(w http.ResponseWriter, r *http.Request) {
	duration := r.URL.Query().Get("duration")
	exchanges := new(domain.Exchange)
	var err error
	if duration == "" {
		exchanges, err = h.use.GetHighestPriceBySym(r.Context(), r.PathValue("symbol"))
	} else {
		timeDuration, err := time.ParseDuration(duration)
		if err != nil {
			http.Error(w, "invalid duration", http.StatusBadRequest)
			return
		}
		exchanges = nil
	}
	w.Header().Set("Content-Type", "application/json")
	if err == nil {
		if err := json.NewEncoder(w).Encode(exchanges); err != nil {
			http.Error(w, "failed to encode error", http.StatusInternalServerError)
		}
		return
	}
	if errors.Is(err, domain.ErrSymbolNotFound) {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err := json.NewEncoder(w).Encode(myError{Error: err.Error()}); err != nil {
		http.Error(w, "failed to encode error", http.StatusInternalServerError)
	}
}
