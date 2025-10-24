package exchange

import (
	"context"
	"encoding/json"
	"log/slog"
	"net"
	"time"

	"marketflow/internal/domain"
)

func (s *stream) Subscribe(ctx context.Context) <-chan *domain.Exchange {
	outCh := make(chan *domain.Exchange)
	go func() {
		defer close(outCh)
		for {
			conn, err := net.DialTimeout("tcp", s.addr, s.interval)
			if err != nil {
				slog.Error("adapter", "exchange", err)
				select {
				case <-ctx.Done():
					return
				case <-time.After(s.interval):
					continue
				}
			}

			dec := json.NewDecoder(conn)
			for {
				select {
				case <-ctx.Done():
					conn.Close()
					return
				default:
					ex := s.get()
					if err := dec.Decode(ex); err != nil {
						conn.Close()
						slog.Error("adapter", "exchange", err)
						break
					}
					ex.Source = s.exName
					outCh <- ex
				}
			}
		}
	}()
	return outCh
}
