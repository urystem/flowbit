package streams

import (
	"context"

	"marketflow/internal/domain"
)

type StreamsInter interface {
	StartStreams(ctx context.Context) error
	StopJustStreams()
	StartTestStream(ctx context.Context) error
	StopTestStream()
	ReturnCh() <-chan *domain.Exchange
	StopStreams()
}
