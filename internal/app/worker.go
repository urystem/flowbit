package app

import (
	"context"
	"log/slog"
	"sync"

	"marketflow/internal/domain"
	"marketflow/internal/ports/outbound"
)

type worker struct {
	jobs <-chan *domain.Exchange
	wg   *sync.WaitGroup
	rdb  outbound.RedisInterForWorkers
	// psql
	quit chan struct{}
}

func (app *workerControl) initWorker(rdb outbound.RedisInterForWorkers, jobs <-chan *domain.Exchange, wg *sync.WaitGroup) *worker {
	return &worker{
		rdb:  rdb,
		jobs: jobs,
		quit: make(chan struct{}),
		wg:   wg,
	}
}

func (w *worker) Start() {
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
				slog.Error(err.Error())
			}
		}
	}
}

func (w *worker) Stop() {
	close(w.quit)
}
