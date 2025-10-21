package one

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"marketflow/internal/services/batcher"
)

func (one *oneMinute) RunWithGo(ctx context.Context, batcher batcher.InsertAndStatus) {
	if one.batcher != nil {
		return
	}
	one.batcher = batcher

	from := time.Now().Truncate(time.Minute)
	next := from.Add(time.Minute) // ближайшая "ровная" минута
	timer := time.NewTimer(time.Until(next))
	defer slog.Info("timer stoped")
	defer timer.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			one.insertAverage(ctx, from, next)
			from = next
			next = next.Add(time.Minute)
			timer.Reset(time.Until(next))
		}
	}
}

func (one *oneMinute) CheckNotAllow() bool {
	return one.notAllow.Load()
}

func (one *oneMinute) WasError() {
	one.wasErr.Store(true)
}

func (one *oneMinute) CollectOldsAndSetAllow(ctx context.Context) error {
	if one.displaced == 0 {
		return fmt.Errorf("%s", "there was not any error")
	}
	to := one.displaced + int(time.Minute)
	ctxR, cancel := context.WithTimeout(ctx, 15*time.Second)
	avgs, err := one.red.GetAllAverages(ctxR, one.displaced, to)
	cancel()
	if err != nil {
		return err
	}
	fromT := time.UnixMilli(int64(one.displaced))
	avgsDB, err := one.db.GetAverageAndDelete(ctx, fromT, time.UnixMilli(int64(to)))
	if err != nil {
		return err
	}
	avgs = one.merger(avgs, avgsDB)

	ctxS, cancel := context.WithTimeout(ctx, 27*time.Second)
	err = one.db.SaveWithCopyFrom(ctxS, avgs, fromT)
	cancel()
	if err != nil {
		return err
	}

	one.notAllow.Store(false)
	one.displaced = 0
	return nil
}
