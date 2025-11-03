package usecase

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"marketflow/internal/domain"
)

func (u *myUsecase) GetLowestPriceBySym(ctx context.Context, sym string) (*domain.Exchange, error) {
	ex, err := u.db.GetLowestPriceBySym(ctx, sym)
	if err == nil {
		return ex, nil
	}
	if errors.Is(err, domain.ErrSymbolNotFound) {
		slog.Info("usecase", "get lowest by symbol in sql average(may be too early request)", err)
	} else {
		slog.Error("usecase", "get lowest by symbol in sql average", err)
		return nil, err
	}
	if !u.one.IsNotWorking() {
		ex, err := u.rdb.GetLowestPriceWithAlign(ctx, 0, sym)
		if err == nil {
			return ex, nil
		}
		if errors.Is(err, domain.ErrSymbolNotFound) {
			return nil, domain.ErrSymbolNotFound
		}
		slog.Error("usecase", "get lowest by symbol in redis", err)
		return nil, domain.ErrInternal
	}

	u.one.PushDone(ctx)
	ex, err = u.db.GetLowestPriceBySymInBackup(ctx, sym)
	if err == nil {
		return ex, nil
	}
	if errors.Is(err, domain.ErrSymbolNotFound) {
		return nil, domain.ErrSymbolNotFound
	}
	slog.Error("usecase", "get lowest by symbol in sql", err)
	return nil, domain.ErrInternal
}

func (u *myUsecase) GetLowestPriceBySymWithDuration(ctx context.Context, sym string, duration time.Duration) (any, error) {
	now := time.Now()
	from := now.Add(-duration)
	if duration <= time.Minute {
		if !u.one.IsNotWorking() {
			ex, err := u.rdb.GetLowestPriceWithAlign(ctx, int(from.UnixMilli()), sym)
			if err != nil {
				if !errors.Is(err, domain.ErrSymbolNotFound) {
					return nil, domain.ErrInternal
				}
				return nil, err
			}
			return ex, err
		}
		slog.Info("search by backup")
		u.one.PushDone(ctx)
		ex, err := u.db.GetLowestPriceBySymWithDuration(ctx, sym, from)
		if err != nil {
			if !errors.Is(err, domain.ErrSymbolNotFound) {
				slog.Error("usecase", "get highest by duration", err)
				return nil, domain.ErrInternal
			}
			return nil, err
		}
		return ex, nil
	}

	// ceil := from.Truncate(time.Minute).Add(time.Minute)
	rounded := from.Round(time.Minute)
	truncated := now.Truncate(time.Minute)
	ex1, err := u.db.GetLowestPriceBySymWithDurationInAverage(ctx, sym, rounded)
	if err != nil {
		if !errors.Is(err, domain.ErrSymbolNotFound) {
			return nil, domain.ErrInternal
		}
		truncated = now.Add(-duration)
	}
	var ex2 *domain.Exchange
	if !u.one.IsNotWorking() {
		ex2, err = u.rdb.GetLowestPriceWithAlign(ctx, int(truncated.UnixMilli()), sym)
		if err != nil {
			if !errors.Is(err, domain.ErrSymbolNotFound) {
				return nil, domain.ErrInternal
			}
		}
	} else {
		u.one.PushDone(ctx)
		ex2, err = u.db.GetLowestPriceBySymWithDuration(ctx, sym, truncated)
		if err != nil {
			if !errors.Is(err, domain.ErrSymbolNotFound) {
				return nil, domain.ErrInternal
			}
		}
	}
	ex := lowestMerge(ex1, ex2)
	if ex == nil {
		return nil, domain.ErrSymbolNotFound
	}
	return &domain.GetExchange{
		Exchange: ex,
		Info:     "rounded:" + now.Sub(rounded).Round(time.Second).String(),
	}, nil
}

func (u *myUsecase) GetLowestPriceByExSym(ctx context.Context, exName, sym string) (*domain.Exchange, error) {
	ex, err := u.db.GetLowestPriceByExSym(ctx, exName, sym)
	if err != nil {
		return ex, nil
	}
	if errors.Is(err, domain.ErrSymbolNotFound) {
		slog.Info("usecase", "get lowest by symbol in sql average(may be too early request)", err)
	} else {
		slog.Error("usecase", "get lowest by symbol in sql average", err)
		return nil, err
	}
	if !u.one.IsNotWorking() {
		ex, err := u.rdb.GetLowestPriceWithEx(ctx, 0, exName, sym)
		if err == nil {
			return ex, nil
		}
		if errors.Is(err, domain.ErrSymbolNotFound) {
			return nil, domain.ErrSymbolNotFound
		}
		slog.Error("usecase", "get lowest by symbol in redis", err)
		return nil, domain.ErrInternal
	}
	u.one.PushDone(ctx)
	ex, err = u.db.GetLowestPriceByExSymInBackup(ctx, exName, sym)
	if err == nil {
		return ex, nil
	}
	if errors.Is(err, domain.ErrSymbolNotFound) {
		return nil, domain.ErrSymbolNotFound
	}
	slog.Error("usecase", "get lowest by symbol in sql", err)
	return nil, domain.ErrInternal
}

func (u *myUsecase) GetLowestPriceByExSymDuration(ctx context.Context, exName, sym string, dur time.Duration) (any, error) {
	now := time.Now()
	from := now.Add(-dur)
	if dur <= time.Minute {
		if !u.one.IsNotWorking() {
			ex, err := u.rdb.GetLowestPriceWithEx(ctx, int(from.UnixMilli()), exName, sym)
			if err != nil {
				if !errors.Is(err, domain.ErrSymbolNotFound) {
					return nil, domain.ErrInternal
				}
				return nil, err
			}
			return ex, err
		}
		u.one.PushDone(ctx)
		ex, err := u.db.GetLowestPriceByExSymWithDuration(ctx, exName, sym, from)
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
	ex1, err := u.db.GetLowestPriceByExSymWithDurationInAverage(ctx, exName, sym, rounded)
	if err != nil {
		if !errors.Is(err, domain.ErrSymbolNotFound) {
			return nil, domain.ErrInternal
		}
		truncated = now.Add(-dur)
	}
	var ex2 *domain.Exchange
	if !u.one.IsNotWorking() {
		ex2, err = u.rdb.GetLowestPriceWithEx(ctx, int(truncated.UnixMilli()), exName, sym)
		if err != nil {
			if !errors.Is(err, domain.ErrSymbolNotFound) {
				return nil, domain.ErrInternal
			}
		}
	} else {
		u.one.PushDone(ctx)
		ex2, err = u.db.GetLowestPriceByExSymWithDuration(ctx, exName, sym, truncated)
		if err != nil {
			if !errors.Is(err, domain.ErrSymbolNotFound) {
				return nil, domain.ErrInternal
			}
		}
	}
	ex := lowestMerge(ex1, ex2)
	if ex == nil {
		return nil, domain.ErrSymbolNotFound
	}
	return &domain.GetExchange{
		Exchange: ex,
		Info:     "rounded:" + now.Sub(rounded).String(),
	}, nil
}

func lowestMerge(ex1, ex2 *domain.Exchange) *domain.Exchange {
	if ex1 == nil {
		return ex2
	} else if ex2 == nil {
		return ex1
	}
	if ex1.Price > ex2.Price {
		return ex2
	} else if ex1.Price < ex2.Price {
		return ex1
	}
	if ex1.Timestamp > ex2.Timestamp {
		return ex1
	}
	return ex2
}
