package generator

import (
	"context"
	"log/slog"
	"net"
	"sync"
	"time"

	"marketflow/internal/domain"
	"marketflow/internal/ports/outbound"
)

type generator struct {
	genName  string
	addr     string
	interval time.Duration
	get      func() *domain.Exchange
	clients  sync.Map
}

func BuildGenerator(exName, addr string, interval time.Duration, get func() *domain.Exchange) outbound.GeneratorInter {
	return &generator{
		genName:  exName,
		addr:     addr,
		interval: interval,
		get:      get,
	}
}

func (g *generator) Check() error {
	ln, err := net.Listen("tcp", g.addr)
	if err != nil {
		return err
	}
	return ln.Close()
}

func (g *generator) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			go g.connCloserAndClear()
			return
		default:
			ln, err := net.Listen("tcp", g.addr)
			if err != nil {
				slog.Error("test generator", "listen", err)
				continue
			}
			g.goFuncAcceptter(ctx, ln)
		}
	}
}

func (g *generator) connCloserAndClear() {
	g.clients.Range(func(key, _ any) bool {
		conn := key.(net.Conn)
		go conn.Close()
		return true
	})
	g.clients.Clear()
}
