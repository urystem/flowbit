package app

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"marketflow/internal/adapters/driver/exchange"
	"marketflow/internal/domain"
	"marketflow/internal/ports/inbound"
)

func (app *myApp) initStream(cfg inbound.SourcesCfg) (inbound.StreamAppInter, error) {
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

type streams struct {
	wg       sync.WaitGroup        // 3
	outCh    chan *domain.Exchange // all
	ctx      context.Context       // for all 3 exchanes
	cancel   context.CancelFunc
	interval time.Duration
	start    chan struct{}
}

func (s *streams) startStream(addr string) {
	defer s.wg.Done()
	for {
		strm, err := exchange.InitStream(addr)
		if err != nil {
			slog.Error("stream", "reconnecting", err)
			time.Sleep(s.interval)
			continue
		}
		<-s.start
		ch, err := strm.Subscribe(s.ctx)
		if err != nil {
			strm.CloseStream()
			slog.Error("stream", "reconnecting", err)
			time.Sleep(s.interval)
			continue
		}
		s.mergeCh(ch) // after closing channel
		strm.CloseStream()
		if s.ctx.Err() != nil {
			break
		}
	}
	slog.Info("stream", "stopped", addr)
}

func (s *streams) Start(ctx context.Context) <-chan *domain.Exchange {
	s.ctx, s.cancel = context.WithCancel(ctx)
	close(s.start)
	return s.outCh
}

func (s *streams) Stop() {
	s.cancel()
	s.wg.Wait()
}

func (s *streams) mergeCh(ch <-chan *domain.Exchange) {
	for ex := range ch {
		s.outCh <- ex
	}
}

