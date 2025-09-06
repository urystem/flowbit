package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"marketflow/internal/adapters/driver/exchange"
	"marketflow/internal/domain"
	"marketflow/internal/ports/inbound"
)

func (app *myApp) initStream(cfg inbound.SourcesCfg) (inbound.StreamAppInter, error) {
	strm := &streams{
		outCh:    make(chan *domain.Exchange),
		start:    make(chan struct{}),
		interval: cfg.GetInterval(),
	}

	for i, addr := range cfg.GetAddresses() {
		strm.wg.Add(1)
		go strm.startStream(i, addr)
	}

	return strm, nil
}

type streams struct {
	wg          sync.WaitGroup        // 3
	outCh       chan *domain.Exchange // all
	ctx         context.Context       // for all 3 exchanes
	cancelCause context.CancelCauseFunc
	interval    time.Duration
	start       chan struct{}
}

func (s *streams) startStream(i int, addr string) {
	defer s.wg.Done()
	for {
		strm, err := exchange.InitStream(addr)
		if err != nil {
			slog.Error(err.Error())
			s.checkErr(err)
			continue
		}

		<-s.start
		ch, err := strm.Subscribe(s.ctx)
		if err != nil {
			s.checkErr(err)
			strm.CloseStream()
			continue
		}
		s.mergeCh(ch)
		s.checkErr(s.ctx.Err())
		strm.CloseStream()
	}
}

func (s *streams) Start(ctx context.Context) <-chan *domain.Exchange {
	s.ctx, s.cancelCause = context.WithCancelCause(ctx)
	close(s.start)
	return s.outCh
}

func (s *streams) Stop() {
	s.cancelCause(fmt.Errorf("%s", "stopinng"))
	s.wg.Wait()
}

func (s *streams) checkErr(err error) {
	if !errors.Is(err, context.Canceled) {
		slog.Error("reconnecting")
		time.Sleep(s.interval)
	}
}

func (s *streams) mergeCh(ch <-chan *domain.Exchange) {
	for ex := range ch {
		s.outCh <- ex
	}
}
