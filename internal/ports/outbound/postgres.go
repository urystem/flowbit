package outbound

import (
	"context"
	"time"

	"marketflow/internal/domain"
)

type PgxInter interface {
	CloseDB()
	PgxForTimer
	PgxFallBack
	PgxCheck
}

type PgxForTimer interface {
	GetAverageAndDelete(ctx context.Context, from, to time.Time) ([]domain.ExchangeAggregation, error)
	SaveWithCopyFrom(ctx context.Context, avgs []domain.ExchangeAggregation, ti time.Time) error
}

type PgxFallBack interface {
	FallBack(ctx context.Context, exs []*domain.Exchange) error
}

type PgxCheck interface {
	CheckHealth(ctx context.Context) error
}
