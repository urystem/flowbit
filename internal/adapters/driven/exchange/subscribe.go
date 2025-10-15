package exchange

import (
	"context"
	"encoding/json"
	"fmt"

	"marketflow/internal/domain"
)

func (s *stream) Subscribe(ctx context.Context) (<-chan *domain.Exchange, error) {
	s.mu.Lock()
	if s.using {
		return nil, fmt.Errorf("%s", "already working")
	}
	s.using = true
	s.mu.Unlock()
	outCh := make(chan *domain.Exchange)
	go func() {
		defer close(outCh)
		dec := json.NewDecoder(s.con)
		for {
			select {
			case <-ctx.Done():
				s.mu.Lock()
				s.using = false
				s.mu.Unlock()
				return
			default:
				ex := s.get()
				if err := dec.Decode(ex); err != nil { //
					// fmt.Println(err)
					// slog.Error("", err)
					return
				}
				ex.Source = s.exName
				outCh <- ex
			}
		}
	}()
	return outCh, nil
}
