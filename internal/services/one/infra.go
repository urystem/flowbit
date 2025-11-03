package one

import (
	"context"
	"sync/atomic"

	"marketflow/internal/domain"
	"marketflow/internal/ports/outbound"
	syncpool "marketflow/internal/services/syncPool"
)

type oneMinute struct {
	ctx            context.Context
	notWorking     atomic.Bool  // real time //for workers
	knowWasErr     atomic.Bool  // access for switch//for batcher
	wasErrInMinute atomic.Bool  // there was error in this last minute //for one minute
	displaced      atomic.Int64 // was old data in redis and db (time)
	red            outbound.RedisForOne
	db             outbound.PgxForTimerAndBatcher
	channel        <-chan *domain.Exchange // fallback
	insertSignals  chan chan<- struct{}
	batch          [][]any
	putter         syncpool.Putter
}

func NewTimerOneMinute(red outbound.RedisForOne, db outbound.PgxForTimerAndBatcher, ch <-chan *domain.Exchange, putter syncpool.Putter) OneMinuteGlobalInter {
	return &oneMinute{
		red:           red, // redis
		db:            db,  // sql
		channel:       ch,  // worker-pool
		insertSignals: make(chan chan<- struct{}),
		putter:        putter,
	}
}

// skipper
