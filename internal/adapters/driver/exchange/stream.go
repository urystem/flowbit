package exchange

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"

	"marketflow/internal/domain"
	"marketflow/internal/ports/inbound"
)

type stream struct {
	exName string
	con    net.Conn
	using  bool
	mu     sync.Mutex
}

func InitStream(addr string) (inbound.StreamAdapterInter, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	// conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", "example.com:1234")
	before, _, _ := strings.Cut(addr, ":")
	return &stream{
		exName: before,
		con:    conn,
	}, nil
}

func (s *stream) CloseStream() error {
	return s.con.Close()
}

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
				ctx.Err()
				s.mu.Lock()
				s.using = false
				s.mu.Unlock()
				return
			default:
				ex := new(domain.Exchange)
				if err := dec.Decode(ex); err != nil { //
					fmt.Println(err)
					return
				}
				ex.Source = s.exName
				outCh <- ex
			}
		}
	}()
	return outCh, nil
}
