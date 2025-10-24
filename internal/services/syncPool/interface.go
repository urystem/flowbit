package syncpool

import "marketflow/internal/domain"

type Getter interface {
	GetNewExchange() *domain.Exchange
}

type Putter interface {
	GetFuncExchange() func(*domain.Exchange)
}

type SyncPoolExchange interface {
	Putter
	Getter
}
