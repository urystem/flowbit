package one

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

func (one *oneMinute) goFuncBatcher() {
	for {
		if len(one.batch) > 511 {
			one.insertBatches()
		}
		select {
		case <-one.ctx.Done():
			return
		case <-one.working:
			fmt.Println("redis worked")
			one.insertBatches()
		case ex := <-one.channel:
			if !one.knowWasErr.Load() {
				one.knowWasErr.Store(true)     // for batcher for this
				one.notWorking.Store(true)     // for workers
				one.wasErrInMinute.Store(true) // fo one minute
				slog.Info("redis not working now")
				go one.tryReconnect()
			}
			one.batch = append(one.batch, ex)
		}
	}
}

func (one *oneMinute) tryReconnect() {
	ti := time.NewTicker(5 * time.Second)
	defer ti.Stop()
	for {
		select {
		case <-one.ctx.Done():
			return
		case <-ti.C:
			err := one.red.CheckHealth(one.ctx)
			if err == nil {
				if one.displaced.Load() != 0 {
					err := one.collectOldsAndSetAllow(one.ctx)
					if err != nil {
						slog.Error("one minute", "try collect old", err)
					}
					slog.Info("old data's average also saved")
					one.wasErrInMinute.Store(true) // tagy 1 ret tekseru ushin,
				}
				one.notWorking.Store(false) // give signal to workers
				time.Sleep(2 * time.Second)
				one.knowWasErr.Store(false) // give signal to batcher
				time.Sleep(2 * time.Second)
				if one.knowWasErr.Load() { // check again
					slog.Error("redis is unstable, try ")
					continue
				}
				one.working <- struct{}{}

				return
			}
		}
	}
}

func (one *oneMinute) insertBatches() error {
	err := one.db.FallBack(one.ctx, one.batch)
	if err != nil {
		slog.Error("fallback", "sql потеря данных", err)
	} else {
		slog.Info("batched")
	}
	putB := one.putter.GetFuncExchange()
	for i := range one.batch {
		putB(one.batch[i])
	}
	one.batch = one.batch[:0]
	return err
}

func (one *oneMinute) collectOldsAndSetAllow(ctx context.Context) error {
	from := one.displaced.Load()
	if from == 0 {
		return fmt.Errorf("%s", "there was not any error")
	}
	to := from + int64(time.Minute)
	avgs, err := one.red.GetAllAverages(ctx, int(from), int(to))
	if err != nil {
		return err
	}
	fromT := time.UnixMilli(one.displaced.Load())
	avgsDB, err := one.db.GetAverageAndDelete(ctx, fromT, time.UnixMilli(int64(to)))
	if err != nil {
		return err
	}
	avgs = one.merger(avgs, avgsDB)

	err = one.db.SaveWithCopyFrom(ctx, avgs, fromT)
	if err != nil {
		return err
	}
	one.displaced.Store(0)
	return nil
}
