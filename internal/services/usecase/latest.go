package usecase

import (
	"context"
	"errors"
	"log/slog"
	"marketflow/internal/domain"
)

// if redis working, we can get latest in redis

/*
if redis not working, we get latest in sql
but we cannot control batch slice, so may be we will get NULL(no any exchanges).
Because batches not fully inserted to sql
*/
func (u *myUsecase) GetLatestBySymbol(ctx context.Context, symbol string) (float64, error) {
	if !u.one.IsNotWorking() {
		price, err := u.rdb.GetLatestPriceBySymbol(ctx, symbol)
		if err == nil {
			return price, nil
		}
		if errors.Is(err, domain.ErrSymbolNotFound) {
			return 0, domain.ErrSymbolNotFound
		}
		slog.Error("usecase", "get latest by symbol in redis", err)
		return 0, domain.ErrInternal
	}

	u.one.PushDone(ctx)

	price, err := u.db.GetLatestPriceBySymbol(ctx, symbol)
	if err == nil {
		return price, nil
	}
	if errors.Is(err, domain.ErrSymbolNotFound) {
		return 0, domain.ErrSymbolNotFound
	}
	slog.Error("usecase", "get latest by symbol in sql", err)
	return 0, domain.ErrInternal
}

func (u *myUsecase) GetLatestPriceByExAndSym(ctx context.Context, ex, sym string) (float64, error) {
	if !u.one.IsNotWorking() {
		price, err := u.rdb.GetLastPriceByExAndSym(ctx, ex, sym)
		if err == nil {
			return price, nil
		}
		if errors.Is(err, domain.ErrSymbolNotFound) {
			return 0, domain.ErrSymbolNotFound
		}
		slog.Error("usecase", "get latest by exchange and symbol in redis", err)
		return 0, domain.ErrInternal
	}

	u.one.PushDone(ctx)
	price, err := u.db.GetLastPriceByExAndSym(ctx, ex, sym)
	if err == nil {
		return price, nil
	}
	if errors.Is(err, domain.ErrSymbolNotFound) {
		return 0, domain.ErrSymbolNotFound
	}
	slog.Error("usecase", "get latest by exchange and symbol in sql", err)
	return 0, domain.ErrInternal
}
