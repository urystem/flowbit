package one

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

func (one *oneMinute) insertAverage(ctx context.Context, from, to time.Time) {
	fromInt, toInt := int(from.UnixMilli()), int(to.UnixMilli())
	ctxR, cancel := context.WithTimeout(ctx, 15*time.Second)
	avgs, err := one.red.GetAllAverages(ctxR, fromInt, toInt)
	cancel()
	if err != nil && one.displaced == 0 {
		one.notAllow.Store(true)
		one.displaced = fromInt
		slog.Error("ticker", "redis error:", err)
		return
	}

	if one.wasErr.Load() || one.displaced != 0 { // wasErr in this minute
		one.wasErr.Store(false)
		err := one.batcher.InsertBatches(ctx)
		if err != nil {
			slog.Error("batch", "insert", err)
		} else {
			avgsDB, err := one.db.GetAverageAndDelete(ctx, from, to)
			if err != nil {
				slog.Error("one", "db erreor", err)
			} else {
				avgs = one.merger(avgs, avgsDB)
			}
		}
	}

	ctx, cancel = context.WithTimeout(ctx, 27*time.Second)
	err = one.db.SaveWithCopyFrom(ctx, avgs, from)
	cancel()
	if err != nil {
		slog.Error("ticker", "psql", err)
	} else {
		fmt.Println(avgs)
		slog.Info("saved to sql")
	}
}
