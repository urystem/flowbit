package usecase

import (
	"context"
	"errors"
	"log/slog"
	"marketflow/internal/domain"
	"time"
)

func (u *myUsecase) GetAveragePriceBySym(ctx context.Context, sym string) (*domain.ExchangeAggregation, error) {
	res, err := u.db.GetAveragePriceBySym(ctx, sym)
	if err == nil {
		return res, nil
	}
	if !errors.Is(err, domain.ErrSymbolNotFound) {
		return nil, domain.ErrInternal
	}

	if u.one.IsNotWorking() {
		res, err = u.rdb.GetAveragePriceWithAlign(ctx, 0, sym)
		if err == nil {
			return res, nil
		}
		if errors.Is(err, domain.ErrSymbolNotFound) {
			return nil, domain.ErrSymbolNotFound
		}
		slog.Error("usecase", "get average by symbol in redis", err)
		return nil, domain.ErrInternal
	}

	u.one.PushDone(ctx)
	res, err = u.db.GetAveragePriceBySymInBackup(ctx, sym)
	if err != nil {
		if errors.Is(err, domain.ErrSymbolNotFound) {
			return nil, err
		}
		slog.Error("usecase", "get average by symbol in sql", err)
		return nil, domain.ErrInternal
	}
	return res, nil
}

func (u *myUsecase) GetAveragePriceBySymWithDuration(ctx context.Context, sym string, dur time.Duration) (*domain.ExchangeAggregation, error) {
	now := time.Now()
	from := now.Add(-dur)
	if dur <= time.Minute {
		if !u.one.IsNotWorking() {
			avg, err := u.rdb.GetAveragePriceWithAlign(ctx, int(from.UnixMilli()), sym)
			if err != nil {
				if !errors.Is(err, domain.ErrSymbolNotFound) {
					return nil, domain.ErrInternal
				}
				return nil, err
			}
			return avg, err
		}
		slog.Info("search by backup for get average")
		u.one.PushDone(ctx)
		ex, err := u.db.GetAveragePriceBySymInBackupTime(ctx, sym, from)
		if err != nil {
			if errors.Is(err, domain.ErrSymbolNotFound) {
				return nil, err
			}
			slog.Error("usecase", "get highest by duration", err)
			return nil, domain.ErrInternal
		}
		return ex, nil
	}

	rounded := from.Round(time.Minute)
	truncated := now.Truncate(time.Minute)
	avg1, err := u.db.GetAveragePriceBySymTime(ctx, sym, truncated)
	if err != nil {
		if !errors.Is(err, domain.ErrSymbolNotFound) {
			return nil, domain.ErrInternal
		}
		truncated = now.Add(-dur)
	}
	var avg2 *domain.ExchangeAggregation
	if !u.one.IsNotWorking() {
		avg2, err = u.rdb.GetAveragePriceWithAlign(ctx, int(truncated.UnixMilli()), sym)
		if err != nil {
			if !errors.Is(err, domain.ErrSymbolNotFound) {
				return nil, domain.ErrInternal
			}
		}
	} else {
		u.one.PushDone(ctx)
		avg2, err = u.db.GetAveragePriceBySymInBackupTime(ctx, sym, truncated)
		if err != nil {
			if !errors.Is(err, domain.ErrSymbolNotFound) {
				return nil, domain.ErrInternal
			}
		}
	}
	avg1 = averageMerge(avg1, avg2)
	if avg1.Count == 0 {
		return nil, domain.ErrSymbolNotFound
	}
	avg1.Timestamp = rounded
	return avg1, nil
}

func (u *myUsecase) GetAveragePriceByExSym(ctx context.Context, exName, sym string) (*domain.ExchangeAggregation, error) {
	agg, err := u.db.GetAveragePriceByExSym(ctx, exName, sym)
	if err == nil {
		return agg, nil
	}
	if errors.Is(err, domain.ErrSymbolNotFound) {
		slog.Info("usecase", "get average by symbol in sql average(may be too early request)", err)
	} else {
		slog.Error("usecase", "get average by symbol in sql average", err)
		return nil, err
	}
	if !u.one.IsNotWorking() {
		aggR, err := u.rdb.GetAveragePriceWithEx(ctx, 0, exName, sym)
		if err == nil {
			return aggR, nil
		}
		if errors.Is(err, domain.ErrSymbolNotFound) {
			return nil, domain.ErrSymbolNotFound
		}
		slog.Error("usecase", "get average by symbol in redis", err)
		return nil, domain.ErrInternal
	}
	u.one.PushDone(ctx)
	aggDB, err := u.db.GetAveragePriceByExSymInBackup(ctx, exName, sym)
	if err == nil {
		return aggDB, nil
	}
	if errors.Is(err, domain.ErrSymbolNotFound) {
		return nil, domain.ErrSymbolNotFound
	}
	slog.Error("usecase", "get average by symbol in sql", err)
	return nil, domain.ErrInternal
}

func (u *myUsecase) GetAveragePriceByExSymDuration(ctx context.Context, exName, sym string, dur time.Duration) (any, error) {
	now := time.Now()
	from := now.Add(-dur)
	if dur <= time.Minute {
		if !u.one.IsNotWorking() {
			ex, err := u.rdb.GetAveragePriceWithEx(ctx, int(from.UnixMilli()), exName, sym)
			if err != nil {
				if !errors.Is(err, domain.ErrSymbolNotFound) {
					return nil, domain.ErrInternal
				}
				return nil, err
			}
			return ex, err
		}
		u.one.PushDone(ctx)
		ex, err := u.db.GetAveragePriceByExSymInBackupTime(ctx, exName, sym, from)
		if err != nil {
			if !errors.Is(err, domain.ErrSymbolNotFound) {
				return nil, domain.ErrInternal
			}
			return nil, err
		}
		return ex, nil
	}
	rounded := from.Round(time.Minute)
	truncated := now.Truncate(time.Minute)
	ex1, err := u.db.GetAveragePriceByExSymTime(ctx, exName, sym, rounded)
	if err != nil {
		if !errors.Is(err, domain.ErrSymbolNotFound) {
			return nil, domain.ErrInternal
		}
		truncated = now.Add(-dur)
	}
	var ex2 *domain.ExchangeAggregation
	if !u.one.IsNotWorking() {
		ex2, err = u.rdb.GetAveragePriceWithEx(ctx, int(truncated.UnixMilli()), exName, sym)
		if err != nil {
			if !errors.Is(err, domain.ErrSymbolNotFound) {
				return nil, domain.ErrInternal
			}
		}
	} else {
		u.one.PushDone(ctx)
		ex2, err = u.db.GetAveragePriceByExSymInBackupTime(ctx, exName, sym, truncated)
		if err != nil {
			if !errors.Is(err, domain.ErrSymbolNotFound) {
				return nil, domain.ErrInternal
			}
		}
	}
	ex := averageMerge(ex1, ex2)
	if ex.Count == 0 {
		return nil, domain.ErrSymbolNotFound
	}
	return ex, nil
}

func averageMerge(aggs ...*domain.ExchangeAggregation) *domain.ExchangeAggregation {
	ans := new(domain.ExchangeAggregation)
	for _, agg := range aggs {
		if agg == nil {
			continue
		}
		ans.AvgPrice = (float64(ans.Count) * ans.AvgPrice) + float64(agg.Count)*agg.AvgPrice/
			float64(ans.Count) + float64(agg.Count)
		ans.Count += ans.Count
		ans.Timestamp = agg.Timestamp
	}
	return ans
}
