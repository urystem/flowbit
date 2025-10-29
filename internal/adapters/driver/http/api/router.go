package router

import "net/http"


func NewRoute(hand any) http.Handler {
	mux := http.NewServeMux()
	// mux.HandleFunc("POST /mode/{status}", hand.Catalog)
	// mux.HandleFunc("GET /health", hand.ServePostImage)

	return mux
	// return mux
}
