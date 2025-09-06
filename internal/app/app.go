package app

import (
	"context"
	"fmt"
	"log/slog"

	"marketflow/internal/ports/inbound"
	"marketflow/internal/ports/outbound"
)

// DI container
type myApp struct {
	ctx            context.Context
	ctxCancelCause context.CancelCauseFunc
	strm           inbound.StreamAppInter // for init //adapter
	red            outbound.RedisInterGlogal
	// ticker inbound.Ticker         // for run
	srv inbound.ServerInter // for init and for run
	// wg  sync.WaitGroup
}

func InitApp(ctx context.Context, cfg inbound.Config) (inbound.AppInter, error) {
	app := &myApp{}
	app.ctx, app.ctxCancelCause = context.WithCancelCause(ctx)

	strm, err := app.initStream(cfg.GetSourcesCfg())
	if err != nil {
		return nil, err
	}
	app.strm = strm

	myRed, err := app.initRedis(ctx, cfg.GetRedisConfig())
	if err != nil {
		return nil, err
	}
	app.red = myRed
	return app, nil
}

func (app *myApp) Shutdown(ctx context.Context) error {
	app.ctxCancelCause(fmt.Errorf("%s", "stopping"))
	app.strm.Stop()
	app.red.CloseRedis()
	return nil
}

func (app *myApp) Run() error {
	slog.Info("server starting")
	ch := app.strm.Start(app.ctx)
	for {
		fmt.Println(<-ch)
	}
	// app.initTicker()
	// time.Sleep(10 * time.Minute)
	return nil
	// return app.srv.ListenServe()
}
