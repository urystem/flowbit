package streams

import (
	"context"
	"sync"
	"time"

	"marketflow/internal/config"
	"marketflow/internal/domain"
)

type streams struct {
	ctx       context.Context // for all 3 exchanes
	cancel    context.CancelFunc
	wg        sync.WaitGroup        // 3
	collectCh chan *domain.Exchange // all
	interval  time.Duration
	start     chan struct{}
	pool      sync.Pool
}

func InitStreams(cfg config.SourcesCfg) (StreamsInter, <-chan *domain.Exchange, error) {
	strm := &streams{
		collectCh: make(chan *domain.Exchange, 128),
		start:     make(chan struct{}),
		interval:  cfg.GetInterval(),
		pool: sync.Pool{
			New: func() any {
				return new(domain.Exchange)
			},
		},
	}

	for _, addr := range cfg.GetAddresses() {
		strm.wg.Add(1)
		// бул жерде sunc.Waitgroup.Go функция болмайды ол тек func() кабылдайды
		go strm.startStream(addr)
	}

	return strm, strm.collectCh, nil
}

func (s *streams) StartStreams(ctx context.Context) {
	s.ctx, s.cancel = context.WithCancel(ctx)
	close(s.start)
}

func (s *streams) StopStreams() {
	s.cancel()
	s.wg.Wait()
}

func (s *streams) get() *domain.Exchange {
	return s.pool.Get().(*domain.Exchange)
}

func (s *streams) ReturnPutFunc() func(*domain.Exchange) {
	return func(ex *domain.Exchange) {
		s.pool.Put(ex)
	}
}
