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
	lastBuckupTime int
}

func InitRickRedis(ctx context.Context, red inbound.RedisConfig) (outbound.RedisInterGlogal, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redisDB:" + red.GetAddr(), // имя сервиса + порт                 // адрес Redis
		Password: red.GetPass(),              // пароль, если есть
		// DB:       0,                          // номер БД (0 по умолчанию)
	})

	return &myRedis{Client: rdb, lastBuckupTime: int(time.Now().UnixMilli())}, rdb.Ping(ctx).Err()
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

// for 60s > target
func (rdb *myRedis) GetByLabel(ctx context.Context, from, to int, keys ...string) ([]domain.Exchange, error) {
	res, err := rdb.TSMRange(
		ctx,
		from,
		to,
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
				Timestamp: myTime,
			})
		}
	}
	return exchanges, nil
}

func (rdb *myRedis) GetAvarages(ctx context.Context) ([]domain.ExchangeAvg, error) {
	mylast := rdb.lastBuckupTime
	rdb.lastBuckupTime = int(time.Now().UnixMilli())  // текущие миллисекунды Unix
	now := rdb.lastBuckupTime - 27
	keys, err := rdb.Keys(ctx, "*").Result()
	if err != nil {
		return nil, err
	}
	bucketDuration := 80000
	avgs := make([]domain.ExchangeAvg, len(keys))
	for i, key := range keys {
		// BucketDuration ар осындай уакыттагы орташа баганы кайтарады
		// res сонда коп ман кайтарады
		// ал маган тек 1 еу керек
		avg, err := rdb.TSRangeWithArgs(ctx, key, mylast, now, &redis.TSRangeOptions{
			Aggregator:     redis.Avg,
			BucketDuration: bucketDuration,
		}).Result()
		if err != nil {
			return nil, err
		} else if len(avg) != 1 {
			fmt.Println(avg)
			return nil, fmt.Errorf("%s%d", "avg aggregation returned not 1 avg, it is", len(avg))
		}

		parts := strings.SplitN(key, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("%s%s", "invalid key:", key)
		}
		avgs[i].Source = parts[0]
		avgs[i].Symbol = parts[1]
		avgs[i].AvgPrice = avg[0].Value
		// avgs[i].AtTime = int64(rdb.lastBuckupTime)
		avgs[i].AtTime = avg[0].Timestamp

		count, err := rdb.TSRangeWithArgs(ctx, key, mylast, now, &redis.TSRangeOptions{
			Aggregator:     redis.Count,
			BucketDuration: bucketDuration,
		}).Result()
		if err != nil {
			return nil, err
		} else if len(count) != 1 {
			return nil, fmt.Errorf("%s%d", "count aggregation returned not 1 avg, it is ", len(count))
		}
		avgs[i].Count = int(count[0].Value)

		minPrice, err := rdb.TSRangeWithArgs(ctx, key, mylast, rdb.lastBuckupTime, &redis.TSRangeOptions{
			Aggregator:     redis.Min,
			BucketDuration: bucketDuration,
		}).Result()
		if err != nil {
			return nil, err
		} else if len(minPrice) != 1 {
			return nil, fmt.Errorf("%s", "min aggregation returned not 1 avg")
		}
		avgs[i].MinPrice = minPrice[0].Value

		maxPrice, err := rdb.TSRangeWithArgs(ctx, key, mylast, rdb.lastBuckupTime, &redis.TSRangeOptions{
			Aggregator:     redis.Max,
			BucketDuration: bucketDuration,
		}).Result()
		if err != nil {
			return nil, err
		} else if len(minPrice) != 1 {
			return nil, fmt.Errorf("%s", "max aggregation returned not 1 avg")
		}
		avgs[i].MaxPrice = maxPrice[0].Value
	}
	rdb.lastBuckupTime++
	return avgs, nil
}

// 12:00
// 12:01
// 12:02
// 12:03 +
// 12:04 +

// 12:01:33 - 12:04:33
