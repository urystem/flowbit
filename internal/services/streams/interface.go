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
	StreamsPutter
	ReturnCh() <-chan *domain.Exchange
}

type StreamsPutter interface {
	ReturnPutFunc() func(*domain.Exchange)
}
