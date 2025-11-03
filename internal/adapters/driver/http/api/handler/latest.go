package handler

import (
	"net/http"
)

type latest interface {
	GetLatestPriceBySymbol(w http.ResponseWriter, r *http.Request)
	GetLatestPriceByExAndSym(w http.ResponseWriter, r *http.Request)
}

func (h *handler) GetLatestPriceBySymbol(w http.ResponseWriter, r *http.Request) {
	ex, err := h.use.GetLatestBySymbol(r.Context(), r.PathValue("symbol"))
	h.writer(w, ex, err)
}

func (h *handler) GetLatestPriceByExAndSym(w http.ResponseWriter, r *http.Request) {
	ex, err := h.use.GetLatestPriceByExAndSym(r.Context(), r.PathValue("exchange"), r.PathValue("symbol"))
	h.writer(w, ex, err)
}
