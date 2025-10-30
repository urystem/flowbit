package api

import (
	"marketflow/internal/adapters/driver/http/api/handler"
	"marketflow/internal/ports/inbound"
	"net/http"
)

func NewRoute(use inbound.UsecaseInter) http.Handler {
	mux := http.NewServeMux()

	hand := handler.NewHandler(use)
	mux.HandleFunc("POST /mode/test", hand.SwitchToTest)
	mux.HandleFunc("POST /mode/live", hand.SwitchToLive)
	// mux.HandleFunc("GET /health", hand.ServePostImage)
	return mux
	// return mux
}
