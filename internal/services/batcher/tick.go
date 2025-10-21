package batcher

import (
	"log/slog"
	"time"
)

func (f *batchCollector) goFunc() {
	defer close(f.channel)
	defer close(f.working)
	for {
		if len(f.batch) > 511 {
			f.InsertBatches(f.ctx)
		}
		select {
		case <-f.ctx.Done():
			return
		case <-f.working:
			f.switcherNotWoring(false)
			go f.workingSleepToDone()
		case ex := <-f.channel:
			if !f.IsNotWorking() {
				f.switcherNotWoring(true)
				slog.Info("redis not working")
				go f.ticker()
			}
			f.batch = append(f.batch, ex)
		}
	}
}

func (f *batchCollector) ticker() {
	ti := time.NewTicker(5 * time.Second)
	defer ti.Stop()
	for {
		select {
		case <-f.ctx.Done():
			return
		case <-ti.C:
			err := f.rdb.CheckHealth(f.ctx)
			if err == nil {
				f.working <- struct{}{}
				return
			}
		}
	}
}

func (f *batchCollector) workingSleepToDone() {
	time.Sleep(5 * time.Second) // wait for channel
	f.InsertBatches(f.ctx)
	slog.Info("redis working and batched to sql of fallback channel")
}

func (f *batchCollector) switcherNotWoring(b bool) {
	f.notWorking.Store(b)
}
