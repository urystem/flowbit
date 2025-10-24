package app

import (
	"context"
	"fmt"
	"log/slog"

	"marketflow/internal/adapters/driven/postgres"
	"marketflow/internal/adapters/driven/redis"
	"marketflow/internal/config"
	"marketflow/internal/ports/inbound"
	"marketflow/internal/ports/outbound"
	"marketflow/internal/services/one"
	"marketflow/internal/services/streams"
	syncpool "marketflow/internal/services/syncPool"
	"marketflow/internal/services/workers"
)

// DI container
type myApp struct {
	ctx            context.Context
	ctxCancelCause context.CancelCauseFunc
	strm           streams.StreamsInter      // service
	red            outbound.RedisInterGlogal // adapter
	workers        workers.WorkerPoolInter   // service
	db             outbound.PgxInter         // adapter
	one            one.OneMinuteGlobalInter  // service
	// srv            inbound.ServerInter // for init and for run
}

func InitApp(ctx context.Context, cfg config.ConfigInter) (inbound.AppInter, error) {
	app := &myApp{}
	getPut := syncpool.NewSyncPoolExchange()
	myStrm, err := app.initStreamsService(cfg.GetSourcesCfg(), getPut)
	if err != nil {
		return nil, err
	}
	app.strm = myStrm

	myRed, err := redis.InitRickRedis(ctx, cfg.GetRedisConfig())
	if err != nil {
		return nil, err
	}
	app.red = myRed

	myDB, err := postgres.InitDB(app.ctx, cfg.GetDBConfig())
	if err != nil {
		return nil, err
	}
	app.db = myDB

	app.workers = workers.InitWorkers(cfg.GetWorkerCfg(), myRed, getPut, myStrm.ReturnCh())

	app.one = one.NewTimerOneMinute(myRed, myDB, app.workers.ReturnChReadOnly(), getPut)
	return app, nil
}

func (app *myApp) Run() error {
	slog.Info("services are starting")
	err := app.one.Run(app.ctx)
	if err != nil {
		return err
	}
	app.workers.Start(app.ctx, app.one)
	app.strm.StartStreams(app.ctx)
	// time.Sleep(10 * time.Second)
	// res, err := app.red.GetByLabel(app.ctx, 0, 0, "exchange=exchange1") // из мапа
	// res, err := app.red.GetAvarages(context.TODO())
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(len(res))
	// 	fmt.Println(res)
	// }
	// app.initTicker()
	// time.Sleep(10 * time.Minute)
	return nil
	// return app.srv.ListenServe()
}

func (app *myApp) Shutdown(ctx context.Context) error {
	app.ctxCancelCause(fmt.Errorf("%s", "stopping"))
	slog.Info("start shutdown")

	app.strm.StopStreams()
	slog.Info("stream")

	app.workers.CleanAll()
	slog.Info("workers")

	app.red.CloseRedis()
	slog.Info("redis")

	app.db.CloseDB()
	slog.Info("db")
	return nil
}
