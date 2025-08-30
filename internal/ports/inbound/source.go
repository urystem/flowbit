package inbound

import "marketflow/internal/domain"

type SourceInter interface {
	Start() (<-chan *domain.Exchange, error)
}
