package workers

import (
	"context"

	"marketflow/internal/domain"
	"marketflow/internal/ports/outbound"
)

type worker struct {
	jobs     <-chan *domain.Exchange
	rdb      outbound.RedisInterForWorkers
	put      func(*domain.Exchange)
	fallback chan<- *domain.Exchange
	quit     chan struct{}
}

func (app *workerControl) initWorker(rdb outbound.RedisInterForWorkers, put func(*domain.Exchange), jobs <-chan *domain.Exchange, fallBack chan<- *domain.Exchange) workerInter {
	return &worker{
		jobs:     jobs,
		rdb:      rdb,
		put:      put,
		quit:     make(chan struct{}),
		fallback: fallBack,
	}
}

func (w *worker) Start() {
	for {
		select {
		case <-w.quit:
			return
		case ex, ok := <-w.jobs:
			if !ok {
				return // канал закрыт
			}
			err := w.rdb.Add(context.TODO(), ex)
			if err != nil {
				w.fallback <- ex
			} else {
				w.put(ex)
			}
		}
	}
}

func (w *worker) Stop() {
	close(w.quit)
}
