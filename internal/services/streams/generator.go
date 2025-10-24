package streams

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"marketflow/internal/domain"
)

func (s *streams) StopTestStream() {
	if s.cancelTest == nil {
		return
	}
	s.cancelTest()
}

func (s *streams) StartTestStream(ctx context.Context) error {
	if s.ctxTest != nil && s.ctxTest.Err() == nil {
		return fmt.Errorf("%s", "test already running")
	} else if s.closedCh.Load() {
		return fmt.Errorf("%s", "channel closed")
	}
	s.ctxTest, s.cancelTest = context.WithCancel(ctx)
	go func() {
		for {
			select {
			case <-s.ctxTest.Done():
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
