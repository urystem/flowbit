package workers

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"marketflow/internal/config"
	"marketflow/internal/domain"
	"marketflow/internal/ports/outbound"
	"marketflow/internal/services/one"
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
	one            one.OneMinuteStatus
	fallCh         chan *domain.Exchange
	closedFallCh   atomic.Bool
}

func InitWorkers(cfg config.WorkerCfg, rdb outbound.RedisInterForWorkers, strm streams.StreamForWorker) WorkerPoolInter {
	return &workerControl{
		maxOrDefWorker: cfg.GetCountOfMaxOrDefWorker(),
		elastic:        cfg.GetBoolElasticWorker(),
		rdb:            rdb,
		strm:           strm,
		interval:       cfg.GetElasticInterval(),
		fallCh:         make(chan *domain.Exchange, 64),
	}
}

func (wc *workerControl) Start(ctx context.Context, rdbStatusAndBatch one.OneMinuteStatus) {
	wc.ctx = ctx
	wc.one = rdbStatusAndBatch
	for range wc.maxOrDefWorker {
		wc.addWorker()
	}
	if wc.elastic {
		go wc.elasTicker()
	}
}

func (wc *workerControl) CleanAll() { // STOP
	if wc.closedFallCh.Load() {
		return
	}
	defer slog.Info("Stop workers ", "count:", len(wc.workers))
	defer wc.wg.Wait()
	defer close(wc.fallCh)
	defer wc.closedFallCh.Store(true)
	for _, w := range wc.workers {
		go w.Stop()
	}
}

func (wc *workerControl) ReturnChReadOnly() <-chan *domain.Exchange {
	return wc.fallCh
}

func (wc *workerControl) giveAsFall(ex *domain.Exchange) {
	if wc.closedFallCh.Load() {
		slog.Error("wtf")
		return
	}
	wc.fallCh <- ex
}

func (wc *workerControl) addWorker() {
	w := wc.initWorker(wc.strm, wc.rdb, wc.one, wc.giveAsFall)
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
			wc.CleanAll()
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
