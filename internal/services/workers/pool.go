package workers

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"marketflow/internal/config"
	"marketflow/internal/domain"
	"marketflow/internal/ports/inbound"
	"marketflow/internal/ports/outbound"
)

type workerControl struct {
	ctx            context.Context
	wg             sync.WaitGroup
	workers        []inbound.Worker
	maxOrDefWorker int
	interval       time.Duration
	elastic        bool
	rdb            outbound.RedisInterForWorkers
	ex             <-chan *domain.Exchange
	fallBack       chan<- *domain.Exchange
}

func InitWorkers(cfg config.WorkerCfg, rdb outbound.RedisInterForWorkers, ex <-chan *domain.Exchange, fallBack chan<- *domain.Exchange) inbound.WorkerPoolInter {
	return &workerControl{
		maxOrDefWorker: cfg.GetCountOfMaxOrDefWorker(),
		elastic:        cfg.GetBoolElasticWorker(),
		rdb:            rdb,
		ex:             ex,
		fallBack:       fallBack,
		interval:       cfg.GetElasticInterval(),
	}
}

func (wc *workerControl) Start(ctx context.Context) {
	wc.ctx = ctx
	for range wc.maxOrDefWorker {
		wc.addWorker()
	}
	if wc.elastic {
		go wc.elasTicker()
	}
}

func (wc *workerControl) addWorker() {
	w := wc.initWorker(wc.rdb, wc.ex, wc.fallBack)
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

func (wc *workerControl) CleanAll() { // STOP
	// for range wc.workers {
	// 	wc.removeWorker()
	// }
	for _, w := range wc.workers {
		go w.Stop()
	}
	wc.wg.Wait()
}
