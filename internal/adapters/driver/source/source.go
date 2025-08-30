package source

import (
	"context"
	"sync"

	"marketflow/internal/domain"
	"marketflow/internal/ports/inbound"
)

type market struct {
	addrs        []string
	countWorkers uint8
	wg           sync.WaitGroup
	outCh        chan *domain.Exchange
	ctx          context.Context
	cancel       context.CancelFunc
}

func InitSource(ctx context.Context, src inbound.SourcesCfg) inbound.SourceInter {
	cctx, cancel := context.WithCancel(ctx)
	addrs := src.GetAddrs()
	count := src.GetCountWorkers()
	return &market{
		addrs:        addrs,
		countWorkers: count,
		outCh:        make(chan *domain.Exchange, len(addrs)*int(count)),
		ctx:          cctx,
		cancel:       cancel,
	}
}

func (m *market) Start() (<-chan *domain.Exchange, error) {
	return m.outCh, nil
}

func (m *market) Stop() {
	m.cancel()
	m.wg.Wait()
	close(m.outCh)
}
