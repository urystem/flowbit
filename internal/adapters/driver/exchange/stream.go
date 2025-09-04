package exchange

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"marketflow/internal/domain"
)

type stream struct {
	con  net.Conn
	used bool
	mu   sync.Mutex
}

func InitStream(addr string) (any, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	// conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", "example.com:1234")

	return &stream{con: conn}, nil
}

func (s *stream) CloseStream() error {
	return s.con.Close()
}

func (s *stream) Subscribe(ctx context.Context) (<-chan *domain.Exchange, error) {
	s.mu.Lock()
	if s.used {
		return nil, fmt.Errorf("%w", "already working")
	}
	s.used = true
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
				s.used = false
				s.mu.Unlock()
				return
			default:
				ex := new(domain.Exchange)
				if err := dec.Decode(ex); err != nil { //
					fmt.Println(err)
					return
				}
				outCh <- ex
			}
		}
	}()
	return outCh, nil
}

// 	dec := json.NewDecoder(conn)

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			conn.Close()
// 			return nil, nil
// 		default:
// 			ex := new(domain.Exchange)
// 			if err = dec.Decode(ex); err != nil {
// 				//
// 				fmt.Println(err)
// 				break
// 			}
// 			outCh <- ex
// 		}
// 	}
// 	conn.Close()
// }

// func (m *market) subscribe(ex *exchange) {
// 	defer m.wg.Done()
// 	for {
// 		select {
// 		case <-m.ctx.Done():
// 			return
// 		default:
// 			conn, err := net.Dial("tcp", net.JoinHostPort(ex.host, fmt.Sprintf("%d", ex.port)))
// 			if err != nil {
// 				log.Fatal("dial error:", err)
// 				continue
// 			}
// 			dec := json.NewDecoder(conn)
// 			for {
// 				ex := new(domain.Exchange)
// 				if err = dec.Decode(ex); err != nil {
// 					//
// 					fmt.Println(err)
// 					break
// 				}
// 				m.outCh <- ex
// 			}

// 			conn.Close()
// 		}
// 	}
// }

// func (e *exchange) worker(dec *json.Decoder, out chan<- *domain.Exchange) {
// 	defer e.wg.Done()
// 	for {
// 		ex := new(domain.Exchange)
// 		if err := dec.Decode(ex); err != nil {
// 			return
// 		}
// 		out <- ex
// 	}
// }
