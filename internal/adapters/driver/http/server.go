package server

import (
	"context"
	"fmt"
	"net/http"

	"marketflow/internal/adapters/driver/http/api"
	"marketflow/internal/config"
	"marketflow/internal/ports/inbound"
	"marketflow/internal/ports/outbound"
)

type server struct {
	*http.Server
}

// type server *httpus.Server

func InitServer(cfg config.ServerCfg, use inbound.UsecaseInter) outbound.ServerInter {
	return &server{&http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.GetPort()),
		Handler: api.NewRoute(use),
	}}
}

func (srv *server) ListenServe() error {
	fmt.Println(srv.Addr)
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
