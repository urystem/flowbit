package handler

import (
	"net/http"
	"time"
)

type average interface {
	GetAveragePriceBySym(w http.ResponseWriter, r *http.Request)
	GetAveragePriceByExSym(w http.ResponseWriter, r *http.Request)
}

func (h *handler) GetAveragePriceBySym(w http.ResponseWriter, r *http.Request) {
	duration := r.URL.Query().Get("period")
	sym := r.PathValue("symbol")
	ctx := r.Context()
	if duration == "" {
		exchange, err := h.use.GetAveragePriceBySym(ctx, sym)
		h.writer(w, exchange, err)
		return
	}
	timeDuration, err := time.ParseDuration(duration)
	if err != nil {
		h.writer(w, nil, err)
		return
	}
	zat, err := h.use.GetAveragePriceBySymWithDuration(ctx, sym, timeDuration)
	h.writer(w, zat, err)
}

func (h *handler) GetAveragePriceByExSym(w http.ResponseWriter, r *http.Request) {
	duration := r.URL.Query().Get("period")
	ex := r.PathValue("exchange")
	sym := r.PathValue("symbol")
	ctx := r.Context()
	if duration == "" {
		exchange, err := h.use.GetAveragePriceByExSym(ctx, ex, sym)
		h.writer(w, exchange, err)
		return
	}
	timeDuration, err := time.ParseDuration(duration)
	if err != nil {
		h.writer(w, nil, err)
		return
	}
	zat, err := h.use.GetAveragePriceByExSymDuration(ctx, ex, sym, timeDuration)
	h.writer(w, zat, err)
}
