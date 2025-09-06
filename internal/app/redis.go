package app

import (
	"context"

	"marketflow/internal/ports/inbound"
	"marketflow/internal/ports/outbound"

	"marketflow/internal/adapters/driven/redis"
)

func (app *myApp) initRedis(ctx context.Context, redisCfg inbound.RedisConfig) (outbound.RedisInterGlogal, error) {
	return redis.InitRickRedis(ctx, redisCfg)
}
