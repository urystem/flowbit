package app

import (
	"context"
	"fmt"
	"log/slog"

	"marketflow/internal/adapters/driven/postgres"
	"marketflow/internal/adapters/driven/redis"
	server "marketflow/internal/adapters/driver/http"
	"marketflow/internal/config"
	"marketflow/internal/ports/inbound"
	"marketflow/internal/ports/outbound"
	"marketflow/internal/services/one"
	"marketflow/internal/services/streams"
	syncpool "marketflow/internal/services/syncPool"
	"marketflow/internal/services/usecase"
	"marketflow/internal/services/workers"
)

// DI container
type myApp struct {
	ctx            context.Context
	ctxCancelCause context.CancelCauseFunc
	logger         *slog.Logger
	strm           streams.StreamsInter      // service
	red            outbound.RedisInterGlogal // adapter
	workers        workers.WorkerPoolInter   // service
	db             outbound.PgxInter         // adapter
	one            one.OneMinuteGlobalInter  // service
	server         outbound.ServerInter      // adapter
}

func InitApp(ctx context.Context, cfg config.ConfigInter, logger *slog.Logger) (inbound.AppInter, error) {
	app := &myApp{
		logger: logger,
	}
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

	myDB, err := postgres.InitDB(ctx, cfg.GetDBConfig())
	if err != nil {
		return nil, err
	}
	app.db = myDB

	app.workers = workers.InitWorkers(cfg.GetWorkerCfg(), myRed, getPut, myStrm.ReturnCh())

	app.one = one.NewTimerOneMinute(myRed, myDB, app.workers.ReturnChReadOnly(), getPut)

	app.server = server.InitServer(cfg.GetServerCfg(), usecase.NewUsecase(myStrm, myDB, myRed, app.one))
	return app, nil
}

func (app *myApp) Run(ctx context.Context) error {
	app.ctx, app.ctxCancelCause = context.WithCancelCause(ctx)
	slog.Info("services are starting")
	app.strm.StartStreams(app.ctx)
	err := app.one.Run(app.ctx)
	if err != nil {
		return err
	}
	app.workers.Start(app.ctx, app.one)

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
	// time.Sleep(1 * time.Minute)
	// app.strm.StopJustStreams()
	// app.strm.StartTestStream()
	// time.Sleep(10 * time.Second)
	// fmt.Println("hiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii")
	// app.strm.StopTestStream()
	// app.strm.StartJustStreams()
	// fmt.Println(err)
	return app.server.ListenServe()
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
