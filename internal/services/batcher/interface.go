package batcher

import "marketflow/internal/domain"

type FallBackInter interface {
	insertBatch
	returnCh
	statusRedisSynced
}

type InsertAndStatus interface {
	insertBatch
	statusRedisSynced
}

type StatusAndFallback interface {
	statusRedisSynced
	returnCh
}

type statusRedisSynced interface {
	IsNotWorking() bool
}

type insertBatch interface {
	InsertBatches() error
}

type returnCh interface {
	GoAndReturnCh() chan<- *domain.Exchange
}
