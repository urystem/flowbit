package handler

import (
	"encoding/json"
	"net/http"
)

type highest interface {
	GetHighestPriceBySym(w http.ResponseWriter, r *http.Request)
}


func (h *handler) GetHighestPriceBySym(w http.ResponseWriter, r *http.Request) {
	price, err := h.use.GetHighestPriceBySym(r.Context(), r.PathValue("symbol"))
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(myError{Error: err.Error()}); err != nil {
			http.Error(w, "failed to encode error", http.StatusInternalServerError)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(&myPrice{Price: price}); err != nil {
		http.Error(w, "failed to encode error", http.StatusInternalServerError)
	}
}
