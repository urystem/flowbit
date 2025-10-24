package syncpool

import (
	"sync"

	"marketflow/internal/domain"
)

type mySyncPool struct {
	sync.Pool
}

func NewSyncPoolExchange() SyncPoolExchange {
	return &mySyncPool{
		sync.Pool{
			New: func() any {
				return new(domain.Exchange)
			},
		},
	}
}

func (p *mySyncPool) GetNewExchange() *domain.Exchange {
	ex, ok := p.Get().(*domain.Exchange)
	if !ok {
		panic("dd")
	}
	return ex
}

func (p *mySyncPool) GetFuncExchange() func(*domain.Exchange) {
	return func(ex *domain.Exchange) {
		p.Put(ex)
	}
}
