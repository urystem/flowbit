package one

import (
	"sync/atomic"

	"marketflow/internal/ports/outbound"
	"marketflow/internal/services/batcher"
)

type oneMinute struct {
	wasErr    atomic.Bool // there was error in this last minute
	displaced int         // time
	notAllow  atomic.Bool // related with displaced int time
	batcher   batcher.InsertAndStatus
	red       outbound.RedisForOne
	db        outbound.PgxForTimer
}

func NewTimerOneMinute(red outbound.RedisForOne, db outbound.PgxForTimer) OneMinuteGlobalInter {
	return &oneMinute{
		red: red,
		db:  db,
	}
}

// skipper
