package workers

import (
	"context"

	"marketflow/internal/domain"
)

type WorkerPoolInter interface {
	Start(ctx context.Context, fall chan<- *domain.Exchange)
	CleanAll()
}

type workerInter interface {
	Start()
	Stop()
}
