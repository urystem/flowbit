package streams

import (
	"context"

	"marketflow/internal/domain"
)

type StreamsInter interface {
	StartStreams(ctx context.Context)
	StopStreams()
	ReturnPutFunc() func(*domain.Exchange)
}
