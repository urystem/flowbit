package inbound

import (
	"context"

	"marketflow/internal/domain"
)

type StreamsInter interface {
	StartStreams(ctx context.Context) <-chan *domain.Exchange
	StopStreams()
}
