package one

import (
	"context"
	"sync/atomic"

	"marketflow/internal/domain"
	"marketflow/internal/ports/outbound"
	"marketflow/internal/services/streams"
)

type oneMinute struct {
	ctx            context.Context
	notWorking     atomic.Bool  // real time
	wasErrInMinute atomic.Bool  // there was error in this last minute
	displaced      atomic.Int64 // was old data in redis and db (time)
	red            outbound.RedisForOne
	db             outbound.PgxForTimerAndBatcher
	channel        <-chan *domain.Exchange // fallback
	batch          []*domain.Exchange
	strm           streams.StreamsPutter
}

func NewTimerOneMinute(red outbound.RedisForOne, db outbound.PgxForTimerAndBatcher, ch <-chan *domain.Exchange, put streams.StreamsPutter) OneMinuteGlobalInter {
	return &oneMinute{
		red:     red, // redis
		db:      db,  // sql
		channel: ch,  // worker-pool
		strm:    put, // stream
	}
}

// skipper
