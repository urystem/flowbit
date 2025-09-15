package app

import (
	"context"
	"fmt"
	"sync"
	"time"

	"marketflow/internal/domain"
	"marketflow/internal/ports/inbound"
	"marketflow/internal/ports/outbound"
)

type workerControl struct {
	ctx            context.Context
	cancelFuncC    context.CancelCauseFunc
	wg             sync.WaitGroup
	workers        []*worker
	maxOrDefWorker int
	interval       time.Duration
	elastic        bool
	rdb            outbound.RedisInterForWorkers
	ex             <-chan *domain.Exchange
}

type WorkerInter interface {
	Start(ctx context.Context)
}

func (app *myApp) initWorkers(cfg inbound.WorkerCfg, rdb outbound.RedisInterForWorkers, ex <-chan *domain.Exchange) WorkerInter {
	return &workerControl{
		maxOrDefWorker: cfg.GetCountOfMaxOrDefWorker(),
		elastic:        cfg.GetBoolElasticWorker(),
		rdb:            rdb,
		ex:             ex,
		interval:       cfg.GetElasticInterval(),
	}
}

func (wc *workerControl) Start(ctx context.Context) {
	wc.ctx, wc.cancelFuncC = context.WithCancelCause(ctx)
	for range wc.maxOrDefWorker {
		wc.addWorker()
	}
	if wc.elastic {
		go wc.elasTicker()
	}
}

func (wc *workerControl) addWorker() {
	w := wc.initWorker(wc.rdb, wc.ex, &wc.wg)
	wc.workers = append(wc.workers, w)
	w.Start()
}

func (wc *workerControl) removeWorker() {
	w := wc.workers[len(wc.workers)-1]
	w.Stop()
	wc.workers = wc.workers[:len(wc.workers)-1]
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
			if kanallen == 0 { // for panic
				continue
			} else if qatynas := float32(wlen) / float32(kanallen); qatynas < 1.1 && wlen < int(wc.maxOrDefWorker) {
				wc.addWorker()
				fmt.Println("⚡ Added worker, total:", len(wc.workers))
			} else if qatynas > 1.5 {
				wc.removeWorker()
				fmt.Println("💤 Removed worker, total:", len(wc.workers))
			}
		}
	}
}

func (wc *workerControl) cleanAll() { // STOP
	wc.cancelFuncC(fmt.Errorf("stop worker"))
	for _, w := range wc.workers {
		go w.Stop()
	}
	wc.wg.Wait()
}
