package streams

import (
	"context"
	"fmt"
	"log/slog"
)

func (s *streams) StopTestStream() {
	if !s.testRunning.Load() || s.closedCh.Load() {
		return
	}
	s.cancelTest()
	s.testRunning.Store(false)
	slog.Info("stopped test mode")
}

func (s *streams) StartTestStream() error {
	if s.closedCh.Load() {
		return fmt.Errorf("%s", "channel closed")
	} else if s.testRunning.Load() {
		return fmt.Errorf("%s", "test already running")
	}
	s.testRunning.Store(true)
	ctx, cancel := context.WithCancel(s.ctxMain)
	s.cancelTest = cancel
	go s.generate.Start(ctx)
	go func() {
		for ex := range s.tester.Subscribe(ctx) {
			s.collectCh <- ex
		}
	}()
	return nil
}
