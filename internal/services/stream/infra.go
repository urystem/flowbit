package stream

import (
	"context"
	"sync"
	"time"

	"marketflow/internal/config"
	"marketflow/internal/domain"
	"marketflow/internal/ports/inbound"
)

type streams struct {
	ctx      context.Context // for all 3 exchanes
	cancel   context.CancelFunc
	wg       sync.WaitGroup        // 3
	outCh    chan *domain.Exchange // all
	interval time.Duration
	start    chan struct{}
}

func InitStream(cfg config.SourcesCfg) (inbound.StreamsInter, error) {
	strm := &streams{
		outCh:    make(chan *domain.Exchange, 123),
		start:    make(chan struct{}),
		interval: cfg.GetInterval(),
	}

	for _, addr := range cfg.GetAddresses() {
		strm.wg.Add(1)
		// бул жерде sunc.Waitgroup.Go функция болмайды ол тек func() кабылдайды
		go strm.startStream(addr)
	}

	return strm, nil
}

func (s *streams) StartStreams(ctx context.Context) <-chan *domain.Exchange {
	s.ctx, s.cancel = context.WithCancel(ctx)
	close(s.start)
	return s.outCh
}

func (s *streams) StopStreams() {
	s.cancel()
	s.wg.Wait()
}
