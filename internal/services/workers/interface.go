package workers

import (
	"context"

	"marketflow/internal/domain"
	"marketflow/internal/services/one"
)

type WorkerPoolInter interface {
	ReturnChReadOnly() <-chan *domain.Exchange
	Start(ctx context.Context, batch one.RedisNotWorking)
	CleanAll()
}

type workerInter interface {
	Start()
	Stop()
}
