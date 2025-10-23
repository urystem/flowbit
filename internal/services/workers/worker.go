package workers

import (
	"context"
	"log/slog"

	"marketflow/internal/domain"
	"marketflow/internal/ports/outbound"
	"marketflow/internal/services/one"
	"marketflow/internal/services/streams"
)

type worker struct {
	strm streams.StreamForWorker
	rdb  outbound.RedisInterForWorkers
	one  one.OneMinuteStatus
	fall func(*domain.Exchange)
	quit chan struct{}
}

func (app *workerControl) initWorker(strm streams.StreamForWorker, rdb outbound.RedisInterForWorkers, one one.OneMinuteStatus, fallWrite func(*domain.Exchange)) workerInter {
	return &worker{
		strm: strm,
		rdb:  rdb,
		one:  one,
		fall: fallWrite,
		quit: make(chan struct{}),
	}
}

func (w *worker) Start() {
	jobs := w.strm.ReturnCh()
	put := w.strm.ReturnPutFunc()
	for {
		select {
		case <-w.quit:
			return
		case ex, ok := <-jobs:
			if !ok {
				return // канал закрыт
			}
			if !w.one.IsNotWorking() { // == working
				err := w.rdb.Add(context.TODO(), ex)
				if err != nil {
					slog.Error("worker", "redis add error:", err)
					w.fall(ex)
				} else {
					put(ex)
				}
			} else {
				w.fall(ex)
			}
		}
	}
}

func (w *worker) Stop() {
	close(w.quit)
}
