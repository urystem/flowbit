package redis

import (
	"context"

	"marketflow/internal/ports/inbound"
	"marketflow/internal/ports/outbound"

	"github.com/redis/go-redis/v9"
)

type myRedis struct {
	*redis.Client
}

func InitRickRedis(ctx context.Context, red inbound.RedisConfig) (outbound.RedisInterGlogal, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redisDB:" + red.GetAddr(), // имя сервиса + порт                 // адрес Redis
		Password: red.GetPass(),              // пароль, если есть
		// DB:       0,                          // номер БД (0 по умолчанию)
	})
	return &myRedis{rdb}, rdb.Ping(ctx).Err()
}

func (rdb *myRedis) CloseRedis() error {
	return rdb.Close()
}

func (rdb *myRedis) CheckHealth(ctx context.Context) error {
	return rdb.Ping(ctx).Err()
}

// for 60s > target

// 12:00
// 12:01
// 12:02
// 12:03 +
// 12:04 +

// 12:01:33 - 12:04:33
