package exchange

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"

	"marketflow/internal/domain"
)

type exchange struct {
	host         string
	port         uint16
	countWorkers uint8
	wg           sync.WaitGroup
}

func (m *market) dialConn(ex *exchange) {
	defer m.wg.Done()
	for {
		select {
		case <-m.ctx.Done():
			return
		default:
			conn, err := net.Dial("tcp", net.JoinHostPort(ex.host, fmt.Sprintf("%d", ex.port)))
			if err != nil {
				log.Fatal("dial error:", err)
				continue
			}
			dec := json.NewDecoder(conn)
			for range ex.countWorkers {
				ex.wg.Add(1)
				go ex.worker(dec, m.outCh)
			}
			ex.wg.Wait()
			conn.Close()
		}
	}
}

func (e *exchange) worker(dec *json.Decoder, out chan<- *domain.Exchange) {
	defer e.wg.Done()
	for {
		ex := new(domain.Exchange)
		if err := dec.Decode(ex); err != nil {
			return
		}
		out <- ex
	}
}
