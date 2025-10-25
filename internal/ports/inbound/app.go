package inbound

import "context"

type AppInter interface {
	Run(context.Context) error
	Shutdown(ctx context.Context) error
}
