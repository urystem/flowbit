package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"marketflow/internal/domain"
)

type myError struct {
	Error string `json:"error"`
}

func (h *handler) writer(w http.ResponseWriter, zat any, err error) {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		if errors.Is(err, domain.ErrSymbolNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else if errors.Is(err, domain.ErrInternal) {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		if err := json.NewEncoder(w).Encode(myError{Error: err.Error()}); err != nil {
			http.Error(w, "failed to encode error", http.StatusInternalServerError)
		}
		return
	}
	if err := json.NewEncoder(w).Encode(zat); err != nil {
		http.Error(w, "failed to encode error", http.StatusInternalServerError)
	}
}
