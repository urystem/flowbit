package streams

import (
	"context"

	"marketflow/internal/domain"
)

type StreamsInter interface {
	StartStreams(ctx context.Context) error
	StopStreams()
	StopJustStreams()
	StopTestStream()
	StartTestStream(ctx context.Context) error
	ReturnCh() <-chan *domain.Exchange
}
