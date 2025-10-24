package exchange

import (
	"time"

	"marketflow/internal/domain"
	"marketflow/internal/ports/outbound"
)

type stream struct {
	exName   string
	addr     string
	interval time.Duration
	get      func() *domain.Exchange
}

func InitStream(exName, addr string, interval time.Duration, get func() *domain.Exchange) outbound.StreamAdapterInter {
	return &stream{
		exName:   exName,
		addr:     addr,
		get:      get,
		interval: interval,
	}
}
