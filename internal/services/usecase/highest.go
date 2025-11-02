package usecase

import (
	"context"
	"errors"
	"log/slog"
	"marketflow/internal/domain"
	"time"
)

func (u *myUsecase) GetHighestPriceBySym(ctx context.Context, sym string) (float64, error) {
	price, err := u.db.GetHighestPriceBySym(ctx, sym)
	if err == nil {
		return price, nil
	}

	if errors.Is(err, domain.ErrSymbolNotFound) {
		slog.Info("usecase", "get highest by symbol in redis(may be too early request)", err)
	} else {
		slog.Error("usecase", "get highest by symbol in redis", err)
	}

	if !u.one.IsNotWorking() {
		to := time.Now().UnixMilli()
		from := to - 60000
		price, err = u.rdb.GetHighestPriceBySymWithAlign(ctx, int(from), int(to), sym)
		if err == nil {
			return price, nil
		}
		if errors.Is(err, domain.ErrSymbolNotFound) {
			return 0, domain.ErrSymbolNotFound
		}
		slog.Error("usecase", "get highest by symbol in redis", err)
		return 0, domain.ErrInternal
	}

	price, err = u.db.GetHighestPriceBySymInBackup(ctx, sym)
	if err == nil {
		return price, nil
	}
	if errors.Is(err, domain.ErrSymbolNotFound) {
		slog.Info("usecase", "get highest by symbol in redis(may be too early request)", err)
		return 0, domain.ErrSymbolNotFound
	}
	slog.Error("usecase", "get highest by symbol in redis", err)
	return 0, domain.ErrInternal
}
