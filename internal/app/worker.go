package app

import (
	"context"
	"sync"

	"marketflow/internal/domain"
	"marketflow/internal/ports/outbound"
)

type worker struct {
	id       int
	jobs     <-chan *domain.Exchange
	wg       *sync.WaitGroup
	rdb      outbound.RedisInterForWorkers
	fallback chan<- *domain.Exchange
	quit     chan struct{}
}

func (app *workerControl) initWorker(rdb outbound.RedisInterForWorkers, jobs <-chan *domain.Exchange, fallBack chan<- *domain.Exchange, wg *sync.WaitGroup) *worker {
	return &worker{
		rdb:      rdb,
		jobs:     jobs,
		quit:     make(chan struct{}),
		wg:       wg,
		fallback: fallBack,
	}
}

func (w *worker) Start(id int) {
	w.id = id
	w.wg.Go(w.myFn)
}

func (w *worker) myFn() {
	for {
		select {
		case <-w.quit:
			// fmt.Printf("Worker %d shutting down\n", w.id)
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
