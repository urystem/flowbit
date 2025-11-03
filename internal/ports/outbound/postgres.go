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
	FallBack(ctx context.Context, rows [][]any) error
}

type PgxForUseCase interface {
	CheckHealth(ctx context.Context) error
	pgxLatest
	pgxHighest
}

type pgxLatest interface {
	GetLatestPriceBySymbol(ctx context.Context, symbol string) (*domain.Exchange, error)
	GetLastPriceByExAndSym(ctx context.Context, ex, sym string) (*domain.Exchange, error)
}

type pgxHighest interface {
	GetHighestPriceBySym(ctx context.Context, sym string) (*domain.Exchange, error)
	GetHighestPriceBySymInBackup(ctx context.Context, sym string) (*domain.Exchange, error)
	GetHighestPriceBySymWithDuration(ctx context.Context, sym string, from int64) (*domain.Exchange, error)
	GetHighestPriceBySymWithDurationInAverage(ctx context.Context, sym string, from time.Time) (*domain.Exchange, error)
}
