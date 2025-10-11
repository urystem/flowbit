package outbound

import (
	"context"

	"marketflow/internal/domain"
)

type PgxInter interface {
	SaveWithCopyFrom(ctx context.Context, avgs []domain.ExchangeAvg) error
}
