package bootstrap

// import (
// 	"context"

// 	"marketflow/internal/ports/inbound"
// 	"marketflow/internal/ports/outbound"

// 	"marketflow/internal/adapters/driven/redis"
// )

// func (app *myApp) initRedis(ctx context.Context, redisCfg inbound.RedisConfig) (outbound.RedisInterLocal, error) {
// 	rdb, err := redis.InitRickRedis(ctx, redisCfg)
// 	if err != nil {
// 		return nil, err
// 	}

// 	app.wg.Add(1)
// 	app.srv.RegisterOnShutDown(func() {
// 		defer app.wg.Done()
// 		rdb.CloseRedis()
// 	})

// 	return rdb, nil
// }
