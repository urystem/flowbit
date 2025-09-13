package app

import (
	"context"
	"fmt"
	"sync"
	"time"

	"marketflow/internal/domain"
	"marketflow/internal/ports/outbound"
)

type workerControl struct {
	ctx            context.Context
	wg             sync.WaitGroup
	workers        []*worker
	maxOrDefWorker int
	elastic        bool
	rdb            outbound.RedisInterForWorkers
	ex             <-chan *domain.Exchange
}

func (app *myApp) initWorkers(ctx context.Context, maxWorkers int, rdb outbound.RedisInterForWorkers, ex <-chan *domain.Exchange) any {
	return &workerControl{
		ctx:            ctx,
		maxOrDefWorker: maxWorkers,
		rdb:            rdb,
		ex:             ex,
	}
}

func (wc *workerControl) start() {
	if wc.elastic {
		wc.addWorker() // add 1 worker
		wc.elasTicker()
	} else {
		for range wc.maxOrDefWorker {
			wc.addWorker()
		}
	}
}

func (wc *workerControl) addWorker() {
	w := wc.initWorker(wc.ex, &wc.wg)
	wc.workers = append(wc.workers, w)
	w.Start()
}

func (wc *workerControl) removeWorker() {
	w := wc.workers[len(wc.workers)-1]
	w.Stop()
	wc.workers = wc.workers[:len(wc.workers)-1]
}

func (wc *workerControl) elasTicker() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		kanallen := len(wc.ex)
		wlen := len(wc.workers)
		if kanallen == 0 { // for panic
			continue
		} else if qatynas := float32(wlen) / float32(kanallen); qatynas < 1.1 {
			wc.addWorker()
			fmt.Println("⚡ Added worker, total:", len(wc.workers))
		} else if qatynas > 1.5 {
			wc.removeWorker()
			fmt.Println("💤 Removed worker, total:", len(wc.workers))
		}
	}
}
