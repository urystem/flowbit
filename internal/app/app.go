package bootstrap

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"marketflow/internal/adapters/driver/http/router"
	"marketflow/internal/adapters/driver/http/server"
	"marketflow/internal/ports/inbound"

	"marketflow/internal/adapters/driver/http/handler"
)

// DI container
type myApp struct {
	ticker inbound.Ticker
	srv    inbound.ServerInter
	wg     sync.WaitGroup
}

func InitApp(ctx context.Context, cfg inbound.Config) (inbound.AppInter, error) {
	// init server
	srvCfg := cfg.GetServerCfg()
	mySrv := server.InitServer(srvCfg)

	// init app
	app := &myApp{srv: mySrv}

	rdbCfg := cfg.GetRedisConfig()
	rdb, err := app.initRedis(ctx, rdbCfg)
	if err != nil {
		return nil, errors.Join(err, app.Shutdown(ctx))
	}
	rdb.SetExchange()

	handler, err := handler.InitHandler(middle, useCase)
	if err != nil {
		return nil, errors.Join(err, app.Shutdown(ctx))
	}

	// init router
	router := router.NewRoute(middle, handler)

	app.srv.SetHandler(router)
	return app, nil
}

func (app *myApp) Shutdown(ctx context.Context) error {
	err := app.srv.ShutdownGracefully(ctx)
	if err == nil {
		app.wg.Wait()
	}
	return err
}

func (app *myApp) Run() error {
	slog.Info("server starting")
	app.initTicker()
	return app.srv.ListenServe()
}
