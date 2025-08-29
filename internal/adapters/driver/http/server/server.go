package server

import (
	"context"
	"fmt"
	"net/http"

	"marketflow/internal/ports/inbound"
)

type server struct {
	*http.Server
}

// type server *http.Server

func InitServer(cfg inbound.ServerCfg) inbound.ServerInter {
	return &server{&http.Server{
		Addr: fmt.Sprintf(":%d", cfg.GetPort()),
	}}
}

func (srv *server) SetHandler(hand http.Handler) {
	srv.Handler = hand
}

func (srv *server) ListenServe() error {
	return srv.ListenAndServe()
}

func (srv *server) ShutdownGracefully(ctx context.Context) error {
	return srv.Shutdown(ctx)
}

func (srv *server) RegisterOnShutDown(f func()) {
	srv.RegisterOnShutdown(f)
}

// func (srv *server) CloseServer() error {
// 	return srv.Close()
// }
