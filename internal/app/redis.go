package app

import (
	"context"

	"marketflow/internal/config"
	"marketflow/internal/ports/outbound"

	"marketflow/internal/adapters/driven/redis"
)

func (app *myApp) initRedis(ctx context.Context, redisCfg config.RedisConfig) (outbound.RedisInterGlogal, error) {
	return redis.InitRickRedis(ctx, redisCfg)
}
