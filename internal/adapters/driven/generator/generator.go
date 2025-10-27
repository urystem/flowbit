package generator

import (
	"context"
	"encoding/json"
	"math/rand"
	"net"
	"time"

	"marketflow/internal/domain"
)

func (g *generator) goFuncAcceptter(ctx context.Context, ln net.Listener) {
	go func() {
		<-ctx.Done()
		ln.Close()
	}()
	go g.pusher(ctx)
	for {
		conn, err := ln.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			continue
		}
		g.clients.Store(conn, struct{}{})
	}
}

func (g *generator) pusher(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(300 * time.Millisecond):
			b, err := json.Marshal(g.generator())
			if err != nil {
				return
			}
			g.writeForAll(b)
		}
	}
}

func (g *generator) generator() *domain.Exchange {
	symbols := []string{
		"BTCUSDT",
		"DOGEUSDT",
		"TONUSDT",
		"SOLUSDT",
		"ETHUSDT",
	}
	ex := g.get()
	ex.Source = "test"
	ex.Symbol = symbols[rand.Intn(len(symbols))]
	ex.Price = 0.5 + rand.Float64()*49999.5
	ex.Timestamp = time.Now().UnixMilli()
	return ex
}

func (g *generator) writeForAll(b []byte) {
	g.clients.Range(func(key, _ any) bool {
		g.writer(key.(net.Conn), b)
		return true
	})
}

func (g *generator) writer(conn net.Conn, b []byte) {
	_, err := conn.Write(b)
	if err != nil {
		go g.clients.Delete(conn)
	}
}
