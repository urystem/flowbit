package stream

import (
	"log/slog"
	"time"

	"marketflow/internal/adapters/driven/exchange"
	"marketflow/internal/domain"
)

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

func (s *streams) mergeCh(ch <-chan *domain.Exchange) {
	for ex := range ch {
		s.outCh <- ex
	}
}
