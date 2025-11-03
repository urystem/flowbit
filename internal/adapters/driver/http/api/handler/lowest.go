package handler

import (
	"net/http"
	"time"
)

type lowest interface {
	GetLowestPriceBySym(w http.ResponseWriter, r *http.Request)
	GetLowestPriceByExSym(w http.ResponseWriter, r *http.Request)
}

func (h *handler) GetLowestPriceBySym(w http.ResponseWriter, r *http.Request) {
	duration := r.URL.Query().Get("period")
	sym := r.PathValue("symbol")
	ctx := r.Context()
	if duration == "" {
		exchange, err := h.use.GetLowestPriceBySym(ctx, sym)
		h.writer(w, exchange, err)
		return
	}
	timeDuration, err := time.ParseDuration(duration)
	if err != nil {
		h.writer(w, nil, err)
		return
	}
	zat, err := h.use.GetLowestPriceBySymWithDuration(ctx, sym, timeDuration)
	h.writer(w, zat, err)
}

func (h *handler) GetLowestPriceByExSym(w http.ResponseWriter, r *http.Request) {
	duration := r.URL.Query().Get("period")
	ex := r.PathValue("exchange")
	sym := r.PathValue("symbol")
	ctx := r.Context()
	if duration == "" {
		exchange, err := h.use.GetLowestPriceByExSym(ctx, ex, sym)
		h.writer(w, exchange, err)
		return
	}
	timeDuration, err := time.ParseDuration(duration)
	if err != nil {
		h.writer(w, nil, err)
		return
	}
	zat, err := h.use.GetLowestPriceByExSymDuration(ctx, ex, sym, timeDuration)
	h.writer(w, zat, err)
}
