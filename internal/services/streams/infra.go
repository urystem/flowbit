package streams

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"marketflow/internal/domain"
	"marketflow/internal/ports/outbound"
	syncpool "marketflow/internal/services/syncPool"
)

type streams struct {
	ctxStrm      context.Context
	cancelStream context.CancelFunc
	ctxTest      context.Context
	cancelTest   context.CancelFunc
	strms        []outbound.StreamAdapterInter
	wg           sync.WaitGroup        // 3
	collectCh    chan *domain.Exchange // all
	closedCh     atomic.Bool
	getter       syncpool.Getter // for generator
}

func InitStreams(strms []outbound.StreamAdapterInter, getter syncpool.Getter) StreamsInter {
	return &streams{
		strms:     strms,
		collectCh: make(chan *domain.Exchange, 64),
		getter:    getter,
	}
}

func (s *streams) StartStreams(ctx context.Context) error {
	if s.ctxStrm != nil && s.ctxStrm.Err() == nil {
		return fmt.Errorf("%s", "already running")
	} else if s.closedCh.Load() {
		return fmt.Errorf("%s", "channel closed")
	}
	s.ctxStrm, s.cancelStream = context.WithCancel(ctx)
	for _, strm := range s.strms {
		s.wg.Add(1)
		go s.mergeCh(strm.Subscribe(s.ctxStrm))
	}
	return nil
}

func (s *streams) mergeCh(ch <-chan *domain.Exchange) {
	defer s.wg.Done()
	for ex := range ch {
		s.collectCh <- ex
	}
}

func (s *streams) StopJustStreams() {
	if s.cancelStream != nil {
		s.cancelStream()
	}
	s.wg.Wait()
}

func (s *streams) StopStreams() {
	s.closedCh.Store(true)
	s.StopJustStreams()
	close(s.collectCh)
}

func (s *streams) ReturnCh() <-chan *domain.Exchange {
	return s.collectCh
}
