package redis

import (
	"context"
	"fmt"
	"time"

	"marketflow/internal/domain"
	"marketflow/internal/ports/inbound"
	"marketflow/internal/ports/outbound"

	"github.com/redis/go-redis/v9"
)

type myRedis struct {
	// ctx context.Context
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

func (rdb *myRedis) Add(ctx context.Context, ex *domain.Exchange) error {
	_, err := rdb.TSAddWithArgs(
		ctx,
		ex.Source+":"+ex.Symbol, // tenge
		ex.Timestamp,            // time
		ex.Price,                // price
		&redis.TSOptions{
			Retention: 70000, // 70 SEC
			// DuplicatePolicy: "LAST",
			Labels: map[string]string{
				"exchange": ex.Source,
				"symbol":   ex.Symbol,
			},
		}, // need to exchange
	).Result()
	return err
}

func (rdb *myRedis) CloseRedis() error {
	return rdb.Close()
}

func (rdb *myRedis) Get(ctx context.Context, keys ...string) ([]domain.Exchange, error) {
	now := time.Now().UnixMilli()
	fmt.Println(now)
	res, err := rdb.TSMRange(
		ctx,
		int(now-60_000),
		int(now),
		keys,
	).Result()
	fmt.Println(err)
	var exchanges []domain.Exchange
	for exSym, serie := range res {
		fmt.Println(exSym, serie, "ddddddddd")
	}
	return exchanges, nil
}
