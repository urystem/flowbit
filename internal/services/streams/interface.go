package streams

import (
	"context"

	"marketflow/internal/domain"
)

type StreamsInter interface {
	ReturnCh() <-chan *domain.Exchange
	StartStreams(ctx context.Context) error
	StopStreams()
	StreamUsecase
}

type StreamUsecase interface {
	StartJustStreams()
	StopJustStreams()
	StartTestStream() error
	StopTestStream()
	CheckHealth() map[string]error
}
