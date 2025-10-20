package workers

import (
	"context"

	"marketflow/internal/ports/outbound"
	"marketflow/internal/services/batcher"
	"marketflow/internal/services/streams"
)

type worker struct {
	strm  streams.StreamForWorker
	rdb   outbound.RedisInterForWorkers
	batch batcher.StatusAndFallback
	quit  chan struct{}
}

func (app *workerControl) initWorker(strm streams.StreamForWorker, rdb outbound.RedisInterForWorkers, batch batcher.StatusAndFallback) workerInter {
	return &worker{
		strm:  strm,
		rdb:   rdb,
		batch: batch,
		quit:  make(chan struct{}),
	}
}

func (w *worker) Start() {
	jobs := w.strm.ReturnCh()
	put := w.strm.ReturnPutFunc()
	fall := w.batch.GoAndReturnCh()
	for {
		select {
		case <-w.quit:
			return
		case ex, ok := <-jobs:
			if !ok {
				return // канал закрыт
			}
			if !w.batch.IsNotWorking() {
				err := w.rdb.Add(context.TODO(), ex)
				if err != nil {
					fall <- ex
				} else {
					put(ex)
				}
			} else {
				fall <- ex
			}
		}
	}
}

func (w *worker) Stop() {
	close(w.quit)
}
