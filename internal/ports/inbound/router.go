package inbound

import "net/http"

type RouteInter interface {
	NewServe() http.Handler
}
