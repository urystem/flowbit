package streams

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"marketflow/internal/domain"
)

func (s *streams) StopTestStream() {
	if s.testRunning.Load() {
		return
	}
	s.cancelTest()
	s.testRunning.Store(false)
}

func (s *streams) StartTestStream() error {
	s.testRunning.Store(true)
	if s.ctxMain == nil {
		return fmt.Errorf("%s", "stream is inactive")
	} else if s.closedCh.Load() {
		return fmt.Errorf("%s", "channel closed")
	}
	ctx, cancel := context.WithCancel(s.ctxMain)
	s.cancelTest = cancel
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				s.collectCh <- s.generator()
			}
		}
	}()
	return nil
}

var symbols = []string{
	"BTCUSDT",
	"DOGEUSDT",
	"TONUSDT",
	"SOLUSDT",
	"ETHUSDT",
}

func (s *streams) generator() *domain.Exchange {
	ex := s.getter.GetNewExchange()
	ex.Source = "test"
	ex.Symbol = symbols[rand.Intn(len(symbols))]
	ex.Price = 0.5 + rand.Float64()*49999.5
	ex.Timestamp = time.Now().UnixMilli()
	return ex
}
