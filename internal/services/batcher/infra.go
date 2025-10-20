package batcher

import (
	"context"
	"sync"

	"marketflow/internal/domain"
	"marketflow/internal/ports/outbound"
)

type batchCollector struct {
	ctx                    context.Context
	sql                    outbound.PgxFallBack
	rdb                    outbound.RedisChecker
	working                chan struct{}
	mutex                  sync.RWMutex
	sendedSignalNotWorking bool
	channel                chan *domain.Exchange // fallback
	batch                  []*domain.Exchange
	put                    func(*domain.Exchange)
}

func NewBatchCollector(ctx context.Context, sql outbound.PgxFallBack, rdb outbound.RedisChecker, put func(*domain.Exchange)) FallBackInter {
	return &batchCollector{
		ctx:     ctx,
		sql:     sql,
		rdb:     rdb,
		channel: make(chan *domain.Exchange),
		batch:   make([]*domain.Exchange, 0, 512),
		working: make(chan struct{}),
		put:     put,
	}
}
