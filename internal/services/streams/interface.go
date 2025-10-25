package streams

import (
	"context"

	"marketflow/internal/domain"
)

type StreamsInter interface {
	ReturnCh() <-chan *domain.Exchange
	StartStreams(ctx context.Context) error
	StreamUsecase
}

type StreamUsecase interface {
	StopJustStreams()
	StartTestStream() error
	StopTestStream()
	StopStreams()
	CheckHealth() map[string]error
}
