package app

import (
	"context"
	"fmt"
	"log/slog"

	"marketflow/internal/adapters/driven/postgres"
	"marketflow/internal/domain"
	"marketflow/internal/ports/inbound"
	"marketflow/internal/ports/outbound"
	"marketflow/internal/services/stream"
	"marketflow/internal/services/workers"

	"marketflow/internal/config"
)

// DI container
type myApp struct {
	ctx            context.Context
	ctxCancelCause context.CancelCauseFunc
	strm           inbound.StreamsInter
	red            outbound.RedisInterGlogal
	workCfg        config.WorkerCfg
	workers        inbound.WorkerPoolInter
	db             outbound.PgxInter
	fallBack       chan *domain.Exchange
	// srv            inbound.ServerInter // for init and for run
}

func InitApp(ctx context.Context, cfg config.ConfigInter) (inbound.AppInter, error) {
	app := &myApp{}
	app.ctx, app.ctxCancelCause = context.WithCancelCause(ctx)

	strm, err := stream.InitStream(cfg.GetSourcesCfg())
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
	app.fallBack = make(chan *domain.Exchange, 512)
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
	return nil
}

func (app *myApp) Run() error {
	slog.Info("server starting")
	uCh := app.strm.StartStreams(app.ctx)

	app.workers = workers.InitWorkers(app.workCfg, app.red, uCh, app.fallBack)
	app.workers.Start(app.ctx)
	go app.timerOneMinute()
	go app.tickerToCheckFallBack()
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
