package exchange

import (
	"context"
	"fmt"
	"sync"

	"marketflow/internal/domain"
	"marketflow/internal/ports/inbound"
)

type market struct {
	exchanges []exchange
	wg        sync.WaitGroup // count of exchanges/addresses/sources  = 3
	outCh     chan *domain.Exchange
	// exPool sync.Pool
	ctx    context.Context
	cancel context.CancelFunc
}

func InitSource(ctx context.Context, cfg inbound.SourcesCfg) inbound.SourceInter {
	cctx, cancel := context.WithCancel(ctx)
	hosts := cfg.GetHosts()
	exchanges := make([]exchange, len(hosts))
	for i, hostName := range hosts {
		exchanges[i] = exchange{
			host:         hostName,
			port:         cfg.GetPort(hostName),
			countWorkers: cfg.GetCountWorkers(hostName),
		}
	}

	return &market{
		exchanges: exchanges,
		outCh:     make(chan *domain.Exchange, cfg.GetCountOfAllWorkers()),
		ctx:       cctx,
		cancel:    cancel,
	}
}

func (m *market) Start() (<-chan *domain.Exchange, error) {
	for i := range m.exchanges {
		m.wg.Add(1)
		go m.dialConn(&m.exchanges[i])
	}
	fmt.Println("ddd")
	return m.outCh, nil
}

func (m *market) Stop() {
	m.cancel()
	m.wg.Wait()
	close(m.outCh)
}
