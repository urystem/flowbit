package outbound

import (
	"context"

	"marketflow/internal/domain"
)

type PgxInter interface {
	SaveWithCopyFrom(ctx context.Context, avg *domain.ExchangeAvg) error
}
