package handler

import (
	"encoding/json"
	"net/http"
)

type latest interface {
	GetLatestPriceBySymbol(w http.ResponseWriter, r *http.Request)
	GetLatestPriceByExAndSym(w http.ResponseWriter, r *http.Request)
}

func (h *handler) GetLatestPriceBySymbol(w http.ResponseWriter, r *http.Request) {
	ex, err := h.use.GetLatestBySymbol(r.Context(), r.PathValue("symbol"))
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(myError{Error: err.Error()}); err != nil {
			http.Error(w, "failed to encode error", http.StatusInternalServerError)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(ex); err != nil {
		http.Error(w, "failed to encode error", http.StatusInternalServerError)
	}
}

func (h *handler) GetLatestPriceByExAndSym(w http.ResponseWriter, r *http.Request) {
	ex, err := h.use.GetLatestPriceByExAndSym(r.Context(), r.PathValue("exchange"), r.PathValue("symbol"))
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(myError{Error: err.Error()}); err != nil {
			http.Error(w, "failed to encode error", http.StatusInternalServerError)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(ex); err != nil {
		http.Error(w, "failed to encode error", http.StatusInternalServerError)
	}
}
