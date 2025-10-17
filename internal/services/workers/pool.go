package workers

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"marketflow/internal/config"
	"marketflow/internal/domain"
	"marketflow/internal/ports/outbound"
)

var allocator = sync.Pool{
	New: func() any {
		return new(domain.Exchange)
	},
}

type workerControl struct {
	ctx            context.Context
	wg             sync.WaitGroup
	workers        []workerInter
	maxOrDefWorker int
	interval       time.Duration
	elastic        bool
	rdb            outbound.RedisInterForWorkers
	put            func(*domain.Exchange)
	ex             <-chan *domain.Exchange
	fallBack       chan<- *domain.Exchange
}

func InitWorkers(cfg config.WorkerCfg, rdb outbound.RedisInterForWorkers, put func(*domain.Exchange), ex <-chan *domain.Exchange) WorkerPoolInter {
	return &workerControl{
		maxOrDefWorker: cfg.GetCountOfMaxOrDefWorker(),
		elastic:        cfg.GetBoolElasticWorker(),
		rdb:            rdb,
		put:            put,
		ex:             ex,
		// fallBack:       fallBack,
		interval: cfg.GetElasticInterval(),
	}
}

func (wc *workerControl) Start(ctx context.Context, fallBack chan<- *domain.Exchange) {
	wc.ctx = ctx
	wc.fallBack = fallBack
	for range wc.maxOrDefWorker {
		wc.addWorker()
	}
	if wc.elastic {
		go wc.elasTicker()
	}
}

func (wc *workerControl) CleanAll() { // STOP
	// for range wc.workers {
	// 	wc.removeWorker()
	// }
	for _, w := range wc.workers {
		go w.Stop()
	}
	wc.wg.Wait()
}

func (wc *workerControl) addWorker() {
	w := wc.initWorker(wc.rdb, wc.put, wc.ex, wc.fallBack)
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
	for {
		select {
		case <-wc.ctx.Done():
			return
		case <-ticker.C:
			kanallen := len(wc.ex)
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
