package one

import (
	"context"

	"marketflow/internal/services/batcher"
)

type OneMinuteGlobalInter interface {
	RunWithGo(ctx context.Context, batcher batcher.InsertAndStatus)
	OneMinuteForBatcher
}

type OneMinuteForBatcher interface {
	CheckNotAllow() bool
	WasError()
	CollectOldsAndSetAllow(ctx context.Context) error
}
