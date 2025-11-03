package usecase

import (
	"context"
	"errors"
	"log/slog"
	"marketflow/internal/domain"
	"time"
)

func (u *myUsecase) GetHighestPriceBySym(ctx context.Context, sym string) (*domain.Exchange, error) {
	ex, err := u.db.GetHighestPriceBySym(ctx, sym)
	if err == nil {
		return ex, nil
	}
	if errors.Is(err, domain.ErrSymbolNotFound) {
		slog.Info("usecase", "get highest by symbol in sql average(may be too early request)", err)
	} else {
		slog.Error("usecase", "get highest by symbol in sql average", err)
	}
	if !u.one.IsNotWorking() {
		ex, err := u.rdb.GetHighestPriceWithAlign(ctx, 0, sym)
		if err == nil {
			return ex, nil
		}
		if errors.Is(err, domain.ErrSymbolNotFound) {
			return nil, domain.ErrSymbolNotFound
		}
		slog.Error("usecase", "get highest by symbol in redis", err)
		return nil, domain.ErrInternal
	}

	u.one.PushDone(ctx)
	ex, err = u.db.GetHighestPriceBySymInBackup(ctx, sym)
	if err == nil {
		return ex, nil
	}
	if errors.Is(err, domain.ErrSymbolNotFound) {
		slog.Info("usecase", "get highest by symbol in sql(may be too early request)", err)
		return nil, domain.ErrSymbolNotFound
	}
	slog.Error("usecase", "get highest by symbol in sql", err)
	return nil, domain.ErrInternal
}

func (u *myUsecase) GetHighestPriceBySymWithDuration(ctx context.Context, sym string, duration time.Duration) (any, error) {
	now := time.Now()
	from := now.Add(-duration)

	if duration <= time.Minute {
		if !u.one.IsNotWorking() {
			ex, err := u.rdb.GetHighestPriceWithAlign(ctx, int(from.UnixMilli()), sym)
			if !errors.Is(err, domain.ErrSymbolNotFound) {
				return nil, domain.ErrInternal
			}
			return ex, err
		}
		ex, err := u.db.GetHighestPriceBySymWithDuration(ctx, sym, from.UnixMilli())
		if !errors.Is(err, domain.ErrSymbolNotFound) {
			return nil, domain.ErrInternal
		}
		return ex, err
	}
	// ceil := from.Truncate(time.Minute).Add(time.Minute)
	rounded := from.Round(time.Minute)
	ex1, err := u.db.GetHighestPriceBySymWithDurationInAverage(ctx, sym, rounded)
	if !errors.Is(err, domain.ErrSymbolNotFound) {
		return nil, domain.ErrInternal
	}
	var ex2 *domain.Exchange
	if !u.one.IsNotWorking() {
		ex2, err = u.rdb.GetHighestPriceWithAlign(ctx, int(now.Truncate(time.Minute).UnixMilli()), sym)
		if err != nil {
			if !errors.Is(err, domain.ErrSymbolNotFound) {
				return nil, domain.ErrInternal
			}
			if ex1 == nil {
				return nil, err
			} else {
				return ex1, nil
			}
		}
	}

	ex := highestMerge(ex1, ex2)
	if ex == nil {
		return nil, domain.ErrInternal
	}
	return &domain.GetExchange{
		Exchange: ex,
		Info:     "rounded:"+time.,
	}, nil
}

func highestMerge(ex1, ex2 *domain.Exchange) *domain.Exchange {
	if ex1 == nil {
		return ex2
	} else if ex2 == nil {
		return ex1
	}
	if ex1.Price > ex2.Price {
		return ex1
	} else if ex1.Price < ex2.Price {
		return ex2
	}
	if ex1.Timestamp > ex2.Timestamp {
		return ex1
	}
	return ex2
}
