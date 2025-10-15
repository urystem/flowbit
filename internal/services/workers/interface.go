package workers

import "context"

type WorkerPoolInter interface {
	Start(ctx context.Context)
	CleanAll()
}

type workerInter interface {
	Start()
	Stop()
}
