package usecase

import (
	"context"
	"errors"
	"log/slog"
	"marketflow/internal/domain"
)

// if redis working, we can get latest in redis
// if redis not working, we get latest in sql
func (u *myUsecase) GetLatestBySymbol(ctx context.Context, symbol string) (*domain.Exchange, error) {
	if !u.one.IsNotWorking() {
		price, err := u.rdb.GetLatestPriceBySymbol(ctx, symbol)
		if err == nil {
			return price, nil
		}
		if errors.Is(err, domain.ErrSymbolNotFound) {
			return nil, domain.ErrSymbolNotFound
		}
		slog.Error("usecase", "get latest by symbol in redis", err)
		return nil, domain.ErrInternal
	}

	u.one.PushDone(ctx)

	price, err := u.db.GetLatestPriceBySymbol(ctx, symbol)
	if err == nil {
		return price, nil
	}
	if errors.Is(err, domain.ErrSymbolNotFound) {
		return nil, domain.ErrSymbolNotFound
	}
	slog.Error("usecase", "get latest by symbol in sql", err)
	return nil, domain.ErrInternal
}

func (u *myUsecase) GetLatestPriceByExAndSym(ctx context.Context, ex, sym string) (*domain.Exchange, error) {
	if !u.one.IsNotWorking() {
		price, err := u.rdb.GetLastPriceByExAndSym(ctx, ex, sym)
		if err == nil {
			return price, nil
		}
		if errors.Is(err, domain.ErrSymbolNotFound) {
			return nil, domain.ErrSymbolNotFound
		}
		slog.Error("usecase", "get latest by exchange and symbol in redis", err)
		return nil, domain.ErrInternal
	}

	u.one.PushDone(ctx)
	price, err := u.db.GetLastPriceByExAndSym(ctx, ex, sym)
	if err == nil {
		return price, nil
	}
	if errors.Is(err, domain.ErrSymbolNotFound) {
		return nil, domain.ErrSymbolNotFound
	}
	slog.Error("usecase", "get latest by exchange and symbol in sql", err)
	return nil, domain.ErrInternal
}
