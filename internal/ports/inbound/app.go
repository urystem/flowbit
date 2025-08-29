package inbound

import "context"

type AppInter interface {
	Run() error
	Shutdown(ctx context.Context) error
}
