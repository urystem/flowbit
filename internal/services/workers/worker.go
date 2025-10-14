package workers

import (
	"context"

	"marketflow/internal/domain"
	"marketflow/internal/ports/inbound"
	"marketflow/internal/ports/outbound"
)

type worker struct {
	jobs     <-chan *domain.Exchange
	rdb      outbound.RedisInterForWorkers
	fallback chan<- *domain.Exchange
	quit     chan struct{}
}

func (app *workerControl) initWorker(rdb outbound.RedisInterForWorkers, jobs <-chan *domain.Exchange, fallBack chan<- *domain.Exchange) inbound.Worker {
	return &worker{
		rdb:      rdb,
		jobs:     jobs,
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
			}
		}
	}
}

func (w *worker) Stop() {
	close(w.quit)
}
