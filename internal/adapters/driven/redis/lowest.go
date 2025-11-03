package redis

import (
	"context"
	"fmt"
	"math"
	"time"

	"marketflow/internal/domain"

	"github.com/redis/go-redis/v9"
)

func (r *myRedis) GetLowestPriceWithAlign(ctx context.Context, from int, sym string) (*domain.Exchange, error) {
	to := int(time.Now().UnixMilli())
	high, err := r.TSMRangeWithArgs(
		ctx,
		from,
		to,
		[]string{"symbol=" + sym},
		&redis.TSMRangeOptions{
			Aggregator:      redis.Min,
			BucketDuration:  to - from,
			Align:           from,
			BucketTimestamp: "end",
			WithLabels:      true,
		}).Result()
	if err != nil {
		return nil, err
	} else if len(high) == 0 {
		return nil, domain.ErrSymbolNotFound
	}

	res := &domain.Exchange{
		Symbol: sym,
		Price:  math.MaxFloat64,
	}
	for _, v := range high {
		ex, err := r.getExchangeName(v[0])
		if err != nil {
			return nil, err
		}

		maxi, ok := v[2].([]any)
		if !ok {
			return nil, fmt.Errorf("%s", "unknown error 0")
		} else if len(maxi) != 1 {
			return nil, fmt.Errorf("%s%v", "no 1 args, it is:", maxi)
		}

		myTime, myPrice, err := r.timeValParser(maxi[0])
		if err != nil {
			return nil, err
		}
		if myPrice < res.Price {
			res.Source = ex
			res.Timestamp = myTime
			res.Price = myPrice
		}
	}
	return res, nil
}

func (r *myRedis) GetLowestPriceWithEx(ctx context.Context, from int, exName, sym string) (*domain.Exchange, error) {
	to := int(time.Now().UnixMilli())
	res, err := r.TSRangeWithArgs(ctx, exName+":"+sym, from, to, &redis.TSRangeOptions{
		Aggregator:      redis.Min,
		BucketDuration:  to - from,
		Align:           from,
		BucketTimestamp: "end",
	}).Result()
	if err != nil {
		return nil, err
	} else if ln := len(res); ln == 0 {
		return nil, domain.ErrSymbolNotFound
	} else if ln > 1 {
		return nil, fmt.Errorf("%s%v", "getLowestWithExSym: more than 1 element", res)
	}
	return &domain.Exchange{
		Source:    exName,
		Symbol:    sym,
		Price:     res[0].Value,
		Timestamp: res[0].Timestamp,
	}, nil
}
