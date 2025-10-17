package fallback

import (
	"context"

	"marketflow/internal/domain"
	"marketflow/internal/ports/outbound"
)

type myFallback struct {
	ctx                    context.Context
	sql                    outbound.PgxFallBack
	rdb                    outbound.RedisChecker
	working                chan struct{}
	sendedSignalNotWorking bool
	channel                chan *domain.Exchange
	batch                  []*domain.Exchange
	put                    func(*domain.Exchange)
}

func NewFallback(ctx context.Context, sql outbound.PgxFallBack, rdb outbound.RedisChecker, put func(*domain.Exchange)) FallBackInter {
	return &myFallback{
		ctx:     ctx,
		sql:     sql,
		rdb:     rdb,
		channel: make(chan *domain.Exchange),
		batch:   make([]*domain.Exchange, 0, 512),
		working: make(chan struct{}),
		put:     put,
	}
}
