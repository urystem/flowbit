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
	ctxMain        context.Context    // like backgroud
	cancelStream   context.CancelFunc // for control for stop
	cancelTest     context.CancelFunc // for control for stop
	testRunning    atomic.Bool
	streamsRunning atomic.Bool
	generate       outbound.GeneratorInter
	tester         outbound.StreamAdapterInter
	strms          []outbound.StreamAdapterInter
	wg             sync.WaitGroup        // 3
	collectCh      chan *domain.Exchange // all
	closedCh       atomic.Bool
	getter         syncpool.Getter // for generator
}

func InitStreams(strms []outbound.StreamAdapterInter, generate outbound.GeneratorInter, test outbound.StreamAdapterInter, getter syncpool.Getter) StreamsInter {
	return &streams{
		strms:     strms,
		generate:  generate,
		tester:    test,
		collectCh: make(chan *domain.Exchange, 64),
		getter:    getter,
	}
}

// for application for once
func (s *streams) StartStreams(ctxMain context.Context) error {
	if s.closedCh.Load() {
		return fmt.Errorf("%s", "channel closed")
	} else if s.ctxMain != nil {
		return fmt.Errorf("%s", "already running")
	}
	s.ctxMain = ctxMain
	s.StartJustStreams()
	return nil
}

func (s *streams) StartJustStreams() {
	if s.closedCh.Load() || s.streamsRunning.Load() {
		return
	}
	s.streamsRunning.Store(true)
	ctx, cancel := context.WithCancel(s.ctxMain)
	s.cancelStream = cancel
	for _, strm := range s.strms {
		s.wg.Add(1)
		go s.mergeCh(strm.Subscribe(ctx))
	}
}

func (s *streams) CheckHealth() map[string]error {
	res := make(map[string]error)
	for _, strm := range s.strms {
		name, err := strm.PingStream()
		res[name] = err
	}
	return res
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
		s.streamsRunning.Store(false)
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
