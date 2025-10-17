package fallback

import (
	"log/slog"
	"time"
)

func (f *myFallback) goFunc() {
	defer close(f.channel)
	defer close(f.working)
	for {
		if len(f.batch) > 511 {
			f.InsertBatches()
		}
		select {
		case <-f.ctx.Done():
			return
		case <-f.working:
			go f.workingSleepToDone()
		case ex := <-f.channel:
			if !f.sendedSignalNotWorking {
				f.sendedSignalNotWorking = true
				slog.Info("redis not working")
				go f.ticker()
			}
			f.batch = append(f.batch, ex)
		}
	}
}

func (f *myFallback) ticker() {
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

func (f *myFallback) workingSleepToDone() {
	time.Sleep(5 * time.Second) // wait for channel
	f.InsertBatches()
	f.sendedSignalNotWorking = false
	slog.Info("redis working")
}
