package outbound

import (
	"context"
	"time"

	"marketflow/internal/domain"
)

type PgxInter interface {
	CloseDB()
	PgxForTimerAndBatcher
	PgxForUseCase
}

type PgxForTimerAndBatcher interface {
	GetAverageAndDelete(ctx context.Context, from, to time.Time) ([]domain.ExchangeAggregation, error)
	SaveWithCopyFrom(ctx context.Context, avgs []domain.ExchangeAggregation, ti time.Time) error
	FallBack(ctx context.Context, exs []*domain.Exchange) error
}

type PgxForUseCase interface {
	CheckHealth(ctx context.Context) error
	pgxLatest
	pgxHighest
}

type pgxLatest interface {
	GetLatestPriceBySymbol(ctx context.Context, symbol string) (float64, error)
	GetLastPriceByExAndSym(ctx context.Context, ex, sym string) (float64, error)
}

type pgxHighest interface {
	GetHighestPriceBySym(ctx context.Context, sym string) (float64, error)
	GetHighestPriceBySymInBackup(ctx context.Context, sym string) (float64, error)

	
}
