package one

import (
	"context"
	"log/slog"
	"time"

	"marketflow/internal/domain"
)

func (one *oneMinute) insertAverage(ctx context.Context, from, to time.Time) {
	fromInt, toInt := int(from.UnixMilli()), int(to.UnixMilli())
	var avgs []domain.ExchangeAggregation
	rdbIsnotWorking := one.IsNotWorking()
	if !rdbIsnotWorking {
		avgsRed, err := one.red.GetAllAverages(ctx, fromInt, toInt)
		if err != nil {
			slog.Error("one minute", "redis error:", err)
		}
		avgs = avgsRed
	} else if one.displaced.Load() == 0 {
		one.displaced.Store(int64(fromInt))
		slog.Error("one minute", "redis not working:", "skipping this minute")
		return
	}
	if one.wasErrInMinute.Load() || one.displaced.Load() != 0 { // wasErr in this minute
		one.wasErrInMinute.Store(false)
		err := one.insertBatches()
		if err != nil {
			slog.Error("one minute", "insert", err)
		} else {
			avgsDB, err := one.db.GetAverageAndDelete(ctx, from, to)
			if err != nil {
				slog.Error("one minute", "db error", err)
			} else if rdbIsnotWorking {
				slog.Info("get average from sql")
				avgs = avgsDB
			} else {
				avgs = one.merger(avgs, avgsDB)
				slog.Info("merged from sql and redis")
			}
		}
	}

	err := one.db.SaveWithCopyFrom(ctx, avgs, from)
	if err != nil {
		slog.Error("one minute", "save average error:", err)
	} else {
		slog.Info("saved to sql")
	}
}

func (one *oneMinute) merger(red, db []domain.ExchangeAggregation) []domain.ExchangeAggregation {
	for i, dv := range db {
		var found bool
		for j, rv := range red {
			if rv.Source == dv.Source && rv.Symbol == dv.Symbol {
				red[j].Count += dv.Count
				red[j].AvgPrice = (rv.AvgPrice*float64(rv.Count) + dv.AvgPrice*float64(dv.Count)) /
					float64(red[j].Count)
				red[j].MinPrice = min(rv.MinPrice, dv.MinPrice)
				red[j].MaxPrice = max(rv.MaxPrice, dv.MaxPrice)
				found = true
				break
			}
		}
		if !found {
			red = append(red, db[i])
		}
	}
	return red
}
