package fallback

import "marketflow/internal/domain"

type FallBackInter interface {
	GoAndReturnCh() chan<- *domain.Exchange
	WithoutCh
}

type WithoutCh interface {
	IsWorking() bool
	InsertBatches() error
}
