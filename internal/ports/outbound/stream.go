package outbound

import (
	"context"

	"marketflow/internal/domain"
)

type StreamAdapterInter interface {
	Subscribe(ctx context.Context) (<-chan *domain.Exchange, error)
	CloseStream() error
}
