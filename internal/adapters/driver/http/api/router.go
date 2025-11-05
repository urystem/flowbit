package api

import (
	"net/http"

	"marketflow/internal/adapters/driver/http/api/handler"
	"marketflow/internal/ports/inbound"
)

func NewRoute(use inbound.UsecaseInter) http.Handler {
	mux := http.NewServeMux()

	hand := handler.NewHandler(use)
	mux.HandleFunc("POST /mode/test", hand.SwitchToTest)
	mux.HandleFunc("POST /mode/live", hand.SwitchToLive)

	mux.HandleFunc("GET /health", hand.CheckHealth)

	mux.HandleFunc("GET /prices/latest/{symbol}", hand.GetLatestPriceBySymbol)
	mux.HandleFunc("GET /prices/latest/{exchange}/{symbol}", hand.GetLatestPriceByExAndSym)

	mux.HandleFunc("GET /prices/highest/{symbol}", hand.GetHighestPriceBySym)
	mux.HandleFunc("GET /prices/highest/{exchange}/{symbol}", hand.GetHighestPriceByExSym)

	mux.HandleFunc("GET /prices/lowest/{symbol}", hand.GetLowestPriceBySym)
	mux.HandleFunc("GET /prices/lowest/{exchange}/{symbol}", hand.GetLowestPriceByExSym)

	mux.HandleFunc("GET /prices/average/{symbol}", hand.GetAveragePriceBySym)
	mux.HandleFunc("GET /prices/average/{exchange}/{symbol}", hand.GetAveragePriceByExSym)
	return mux
}

