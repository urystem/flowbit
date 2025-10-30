package outbound

import "context"

type ServerInter interface {
	ListenServe() error
	ShutdownGracefully(ctx context.Context) error
	RegisterOnShutDown(f func())
	// CloseServer() error
}
