package outbound

import (
	"context"
	"time"

	"marketflow/internal/domain"
)

type PgxInter interface {
	CloseDB()
	PgxForTimerAndBatcher
	PgxCheck
	PgxForUseCase
}

type PgxForTimerAndBatcher interface {
	GetAverageAndDelete(ctx context.Context, from, to time.Time) ([]domain.ExchangeAggregation, error)
	SaveWithCopyFrom(ctx context.Context, avgs []domain.ExchangeAggregation, ti time.Time) error
	FallBack(ctx context.Context, exs []*domain.Exchange) error
	PgxCheck
}

type PgxCheck interface {
	CheckHealth(ctx context.Context) error
}

type PgxForUseCase interface {
	PgxCheck
	GetLatestPriceBySymbol(ctx context.Context, symbol string) (float64, error)
	GetLastPriceByExAndSym(ctx context.Context, ex, sym string) (float64, error)
}
