package redis

import (
	"context"
	"fmt"
	"marketflow/internal/domain"

	"github.com/redis/go-redis/v9"
)

// TS.MRANGE 1761981278848 1761981338848 AGGREGATION MAX 60000 ALIGN 1761981278848 FILTER symbol=BTCUSDT
func (r *myRedis) GetHighestPriceBySymWithAlign(ctx context.Context, from, to int, sym string) (float64, error) {
	high, err := r.TSMRangeWithArgs(
		ctx,
		from,
		to,
		[]string{"symbol=" + sym},
		&redis.TSMRangeOptions{
			Aggregator:     redis.Max,
			BucketDuration: to - from,
			Align:          from,
		}).Result()
	if err != nil {
		return 0, err
	} else if len(high) == 0 {
		return 0, domain.ErrSymbolNotFound
	}
	d := high["exchange1:"+sym]
	arr := d[2].([]any)
	fmt.Println(d[0])
	fmt.Println(d[1])
	fmt.Println(len(arr))
	fmt.Println(arr[0])
	return 0, nil
}
