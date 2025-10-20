package workers

import (
	"context"

	"marketflow/internal/services/batcher"
)

type WorkerPoolInter interface {
	Start(ctx context.Context, batch batcher.StatusAndFallback)
	CleanAll()
}

type workerInter interface {
	Start()
	Stop()
}
