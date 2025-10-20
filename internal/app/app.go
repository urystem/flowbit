package app

import (
	"context"
	"fmt"
	"log/slog"

	"marketflow/internal/adapters/driven/postgres"
	"marketflow/internal/adapters/driven/redis"
	"marketflow/internal/ports/inbound"
	"marketflow/internal/ports/outbound"
	"marketflow/internal/services/batcher"
	"marketflow/internal/services/streams"
	"marketflow/internal/services/workers"

	"marketflow/internal/config"
)

// DI container
type myApp struct {
	ctx            context.Context
	ctxCancelCause context.CancelCauseFunc
	strm           streams.StreamsInter // service
	red            outbound.RedisInterGlogal
	workers        workers.WorkerPoolInter // service
	db             outbound.PgxInter
	// srv            inbound.ServerInter // for init and for run
}

func InitApp(ctx context.Context, cfg config.ConfigInter) (inbound.AppInter, error) {
	app := &myApp{}
	app.ctx, app.ctxCancelCause = context.WithCancelCause(ctx)
	strm, err := streams.InitStreams(cfg.GetSourcesCfg())
	if err != nil {
		return nil, err
	}

	app.strm = strm

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
	// fallBack := fallback.Fallback(myDB, myRed, strm.ReturnPutFunc())
	app.workers = workers.InitWorkers(cfg.GetWorkerCfg(), app.red, strm)
	return app, nil
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

func (app *myApp) Run() error {
	slog.Info("server starting")
	fallB := batcher.NewBatchCollector(app.ctx, app.db, app.red, app.strm.ReturnPutFunc())
	app.strm.StartStreams(app.ctx)
	app.workers.Start(app.ctx, fallB)
	go app.timerOneMinute(fallB)
	// go app.tickerToCheckFallBack()
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
