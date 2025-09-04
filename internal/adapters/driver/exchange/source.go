package exchange

// type markets struct {
// 	exchanges []exchange
// 	chs       []chan *domain.Exchange
// 	wg        sync.WaitGroup // count of exchanges/addresses/sources  = 3
// 	outCh     chan *domain.Exchange
// 	// exPool sync.Pool
// 	ctx    context.Context
// 	cancel context.CancelFunc
// }

// func InitSource(ctx context.Context, cfg inbound.SourcesCfg) inbound.SourceInter {
// 	cctx, cancel := context.WithCancel(ctx)
// 	hosts := cfg.GetHosts()
// 	exchanges := make([]exchange, len(hosts))
// 	for i, hostName := range hosts {
// 		exchanges[i] = exchange{
// 			host:         hostName,
// 			port:         cfg.GetPort(hostName),
// 			countWorkers: cfg.GetCountWorkers(hostName),
// 		}
// 	}

// 	return &market{
// 		exchanges: exchanges,
// 		outCh:     make(chan *domain.Exchange, cfg.GetCountOfAllWorkers()),
// 		ctx:       cctx,
// 		cancel:    cancel,
// 	}
// }

// func (m *market) Start() (<-chan *domain.Exchange, error) {
// 	for i := range m.exchanges {
// 		m.wg.Add(1)
// 		go m.subscribe(&m.exchanges[i])
// 	}
// 	return m.outCh, nil
// }

// func (m *market) Stop() {
// 	m.cancel()
// 	m.wg.Wait()
// 	close(m.outCh)
// }
