package workers

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"marketflow/internal/config"
	"marketflow/internal/ports/outbound"
	"marketflow/internal/services/batcher"
	"marketflow/internal/services/streams"
)

type workerControl struct {
	ctx            context.Context
	wg             sync.WaitGroup
	workers        []workerInter
	maxOrDefWorker int
	interval       time.Duration
	elastic        bool
	rdb            outbound.RedisInterForWorkers
	strm           streams.StreamForWorker
	batch          batcher.StatusAndFallback
}

func InitWorkers(cfg config.WorkerCfg, rdb outbound.RedisInterForWorkers, strm streams.StreamForWorker) WorkerPoolInter {
	return &workerControl{
		maxOrDefWorker: cfg.GetCountOfMaxOrDefWorker(),
		elastic:        cfg.GetBoolElasticWorker(),
		rdb:            rdb,
		strm:           strm,
		interval:       cfg.GetElasticInterval(),
	}
}

func (wc *workerControl) Start(ctx context.Context, rdbStatusAndBatch batcher.StatusAndFallback) {
	wc.ctx = ctx
	wc.batch = rdbStatusAndBatch
	for range wc.maxOrDefWorker {
		wc.addWorker()
	}
	if wc.elastic {
		go wc.elasTicker()
	}
}

func (wc *workerControl) CleanAll() { // STOP
	defer slog.Info("Stop workers ", "count:", len(wc.workers))
	defer wc.wg.Wait()
	for _, w := range wc.workers {
		go w.Stop()
	}
}

func (wc *workerControl) addWorker() {
	w := wc.initWorker(wc.strm, wc.rdb, wc.batch)
	wc.workers = append(wc.workers, w)
	wc.wg.Go(w.Start)
	slog.Info("⚡ Added worker", "total:", len(wc.workers))
}

func (wc *workerControl) removeWorker() {
	w := wc.workers[len(wc.workers)-1]
	w.Stop()
	wc.workers = wc.workers[:len(wc.workers)-1]
	slog.Info("💤 Removed worker", "total:", len(wc.workers))
}

func (wc *workerControl) elasTicker() {
	ticker := time.NewTicker(wc.interval)
	defer ticker.Stop()
	ch := wc.strm.ReturnCh()
	for {
		select {
		case <-wc.ctx.Done():
			return
		case <-ticker.C:
			kanallen := len(ch)
			wlen := len(wc.workers)
			if kanallen == 0 && wlen > 1 { // for panic
				wc.removeWorker()
			} else if qatynas := float32(wlen) / float32(kanallen); qatynas < 1.1 && wlen < int(wc.maxOrDefWorker) {
				wc.addWorker()
			} else if qatynas > 1.5 {
				wc.removeWorker()
			}
		}
	}
}
