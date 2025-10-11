package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"marketflow/internal/adapters/driven/postgres"
	"marketflow/internal/ports/inbound"
	"marketflow/internal/ports/outbound"
)

// DI container
type myApp struct {
	ctx            context.Context
	ctxCancelCause context.CancelCauseFunc
	strm           inbound.StreamAppInter // for init //adapter
	red            outbound.RedisInterGlogal
	workCfg        inbound.WorkerCfg
	workers        WorkerInter
	db             outbound.PgxInter
	// srv            inbound.ServerInter // for init and for run
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
	app.workCfg = cfg.GetWorkerCfg()

	myDB, err := postgres.InitDB(app.ctx, cfg.GetDBConfig())
	if err != nil {
		return nil, err
	}
	app.db = myDB
	return app, nil
}

func (app *myApp) Shutdown(ctx context.Context) error {
	app.ctxCancelCause(fmt.Errorf("%s", "stopping"))
	slog.Info("start shutdown")
	app.strm.Stop()
	slog.Info("stream")
	app.workers.CleanAll()
	slog.Info("workers")
	app.red.CloseRedis()
	slog.Info("redis")
	return nil
}

func (app *myApp) Run() error {
	slog.Info("server starting")
	uCh := app.strm.Start(app.ctx)

	app.workers = app.initWorkers(app.workCfg, app.red, uCh)
	app.workers.Start(app.ctx)
	time.Sleep(10 * time.Second)
	// res, err := app.red.GetByLabel(app.ctx, 0, 0, "exchange=exchange1") // из мапа
	// res, err := app.red.GetAvarages(context.TODO())
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(len(res))
	// 	fmt.Println(res)
	// }
	// app.initTicker()
	go app.tickerOneMinute()
	// time.Sleep(10 * time.Minute)
	return nil
	// return app.srv.ListenServe()
}
