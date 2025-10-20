package streams

import (
	"context"

	"marketflow/internal/domain"
)

type StreamsInter interface {
	StartStreams(ctx context.Context)
	StopStreams()
	StreamForWorker
}

type StreamForWorker interface {
	ReturnPutFunc() func(*domain.Exchange)
	ReturnCh() <-chan *domain.Exchange
}
