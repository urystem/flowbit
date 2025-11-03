package redis

import (
	"context"
	"fmt"
	"marketflow/internal/domain"
	"time"

	"github.com/redis/go-redis/v9"
)

// TS.MRANGE 1761981278848 1761981338848 AGGREGATION MAX 60000 ALIGN 1761981278848 FILTER symbol=BTCUSDT
func (r *myRedis) GetHighestPriceWithAlign(ctx context.Context, from int, sym string) (*domain.Exchange, error) {
	to := int(time.Now().UnixMilli())
	high, err := r.TSMRangeWithArgs(
		ctx,
		from,
		to,
		[]string{"symbol=" + sym},
		&redis.TSMRangeOptions{
			Aggregator:      redis.Max,
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
		if myPrice > res.Price {
			res.Source = ex
			res.Timestamp = myTime
			res.Price = myPrice
		}
	}
	return res, nil
}

