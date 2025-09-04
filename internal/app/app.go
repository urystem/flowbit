package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"marketflow/internal/adapters/driver/exchange"
	"marketflow/internal/ports/inbound"
)

// DI container
type myApp struct {
	src    inbound.SourceInter // for init
	ticker inbound.Ticker      // for run
	srv    inbound.ServerInter // for init and for run
	wg     sync.WaitGroup
}

func InitApp(ctx context.Context, cfg inbound.Config) (inbound.AppInter, error) {
	// init server
	// srvCfg := cfg.GetServerCfg()
	// mySrv := server.InitServer(srvCfg)

	// init app
	app := &myApp{}

	// rdbCfg := cfg.GetRedisConfig()
	// rdb, err := app.initRedis(ctx, rdbCfg)
	// if err != nil {
	// 	return nil, errors.Join(err, app.Shutdown(ctx))
	// }
	source := exchange.InitSource(ctx, cfg.GetSourcesCfg())
	app.src = source
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
	ch, err := app.src.Start()
	if err != nil {
		return err
	}
	for {
		fmt.Println(<-ch)
	}
	// app.initTicker()
	time.Sleep(10 * time.Minute)
	return nil
	// return app.srv.ListenServe()
}
