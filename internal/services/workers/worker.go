package workers

import (
	"context"
	"log/slog"

	"marketflow/internal/domain"
	"marketflow/internal/ports/outbound"
	"marketflow/internal/services/one"
	syncpool "marketflow/internal/services/syncPool"
)

type worker struct {
	job    <-chan *domain.Exchange
	rdb    outbound.RedisInterForWorkers
	one    one.OneMinuteStatus
	putter syncpool.Putter
	fall   func(*domain.Exchange)
	quit   chan struct{}
}

func (app *workerControl) initWorker(job <-chan *domain.Exchange, rdb outbound.RedisInterForWorkers, one one.OneMinuteStatus, putter syncpool.Putter, fallWrite func(*domain.Exchange)) workerInter {
	return &worker{
		job:    job,
		rdb:    rdb,
		one:    one,
		putter: putter,
		fall:   fallWrite,
		quit:   make(chan struct{}),
	}
}

func (w *worker) Start() {
	put := w.putter.GetFuncExchange()
	for {
		select {
		case <-w.quit:
			return
		case ex, ok := <-w.job:
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
