package bootstrap

import (
	"context"

	"marketflow/internal/adapters/driven/postgres"
	"marketflow/internal/ports/inbound"
)

func (app *myApp) initSource(ctx context.Context, cfg inbound.SourcesCfg) (inbound.UseCase, error) {
	db, err := postgres.InitDB(ctx, dbCfg)
	if err != nil {
		return nil, err
	}

	app.wg.Add(1)
	app.srv.RegisterOnShutDown(func() {
		defer app.wg.Done()
		db.CloseDB()
	})

	minio, err := minio.InitMinio(ctx, s3Cfg)
	if err != nil {
		return nil, err
	}
	return usecase.InitUsecase(db, minio, session), nil
}
