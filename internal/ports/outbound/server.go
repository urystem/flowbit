package outbound

import (
	"context"
	"net/http"
)

type ServerInter interface {
	SetHandler(hand http.Handler)
	ListenServe() error
	ShutdownGracefully(ctx context.Context) error
	RegisterOnShutDown(f func())
	// CloseServer() error
}
