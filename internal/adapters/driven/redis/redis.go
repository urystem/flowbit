package redis

import (
	"context"
	"fmt"
	"strings"
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
			Retention:       70000, // 70 SEC
			DuplicatePolicy: "LAST",
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
	if err != nil {
		return nil, err
	}
	var exchanges []domain.Exchange
	// serie это 3 интерфейс [3]interface
	// и нам нужна только последный
	for exSym, serie := range res {
		points, ok := serie[2].([]any)
		if !ok {
			return nil, fmt.Errorf("%s", "series is not [3]any")
		}
		myexSym := strings.Split(exSym, ":")
		if len(myexSym) != 2 {
			return nil, fmt.Errorf("%s", "exchange and symbol key not valid")
		}
		for _, p := range points {
			arr, ok := p.([]any)
			if !ok || len(arr) != 2 {
				return nil, fmt.Errorf("%s", "here is not time and price")
			}
			myTime, ok := arr[0].(int64)
			if !ok {
				return nil, fmt.Errorf("%s", "invalid time")
			}
			myPrice, ok := arr[1].(float64)
			if !ok {
				return nil, fmt.Errorf("%s", "invalid price")
			}
			exchanges = append(exchanges, domain.Exchange{
				Source:    myexSym[0],
				Symbol:    myexSym[1],
				Price:     myPrice,
				Timestamp: uint64(myTime),
			})
		}
	}
	return exchanges, nil
}
