package inbound

import "context"

type WorkerPoolInter interface {
	Start(ctx context.Context)
	CleanAll()
}

type Worker interface {
	Start()
	Stop()
}
