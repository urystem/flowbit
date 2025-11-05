package redis

import (
	"context"
	"fmt"
	"marketflow/internal/domain"
	"time"

	"github.com/redis/go-redis/v9"
)

func (r *myRedis) GetAveragePriceWithAlign(ctx context.Context, from int, sym string) (*domain.ExchangeAggregation, error) {
	to := int(time.Now().UnixMilli())
	averages, err := r.TSMRangeWithArgs(
		ctx,
		from,
		to,
		[]string{"symbol=" + sym},
		&redis.TSMRangeOptions{
			Aggregator:      redis.Avg,
			BucketDuration:  to - from,
			Align:           from,
			BucketTimestamp: "end",
			WithLabels:      true,
		}).Result()
	if err != nil {
		return nil, err
	} else if len(averages) == 0 {
		return nil, domain.ErrSymbolNotFound
	}

	res := &domain.ExchangeAggregation{
		Symbol: sym,
	}

	for exSym, v := range averages {
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
		count, err := r.getCount(ctx, exSym, from, to)
		if err != nil {
			return nil, err
		}
		res.AvgPrice = ((float64(res.Count) * res.AvgPrice) + float64(count)*myPrice) /
			(float64(res.Count) + float64(count))
		res.Count += count
		res.Timestamp = time.UnixMilli(myTime)
		if res.Source == "" {
			res.Source += ex
		} else {
			res.Source += ", " + ex
		}
	}
	return res, nil
}

func (r *myRedis) GetAveragePriceWithEx(ctx context.Context, from int, exName, sym string) (*domain.ExchangeAggregation, error) {
	to := int(time.Now().UnixMilli())
	average, err := r.TSRangeWithArgs(
		ctx,
		exName+":"+sym,
		from,
		to,
		&redis.TSRangeOptions{
			Aggregator:      redis.Avg,
			BucketDuration:  to - from,
			Align:           from,
			BucketTimestamp: "end",
		},
	).Result()
	if err != nil {
		return nil, err
	} else if len(average) == 0 {
		return nil, domain.ErrSymbolNotFound
	}

	return &domain.ExchangeAggregation{
		Source:    exName,
		Symbol:    sym,
		AvgPrice:  average[0].Value,
		Timestamp: time.UnixMilli(average[0].Timestamp),
	}, nil
}
