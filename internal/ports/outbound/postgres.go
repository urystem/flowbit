package outbound

import (
	"context"
	"time"

	"marketflow/internal/domain"
)

type PgxInter interface {
	CheckHealth(ctx context.Context) error
	CloseDB()
	SaveWithCopyFrom(ctx context.Context, avgs []domain.ExchangeAggregation, ti time.Time) error
	FallBack(ctx context.Context, exs []domain.Exchange) error
}
