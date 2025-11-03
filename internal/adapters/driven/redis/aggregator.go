package redis

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

// BucketDuration ар осындай уакыттагы орташа баганы кайтарады
// res сонда коп ман кайтарады
// ал маган тек 1 еу керек
func (rdb *myRedis) GetAveragePrice(ctx context.Context, key string, from, to int) (float64, error) {
	avg, err := rdb.TSRangeWithArgs(ctx, key, from, to, &redis.TSRangeOptions{
		Aggregator:     redis.Avg,
		BucketDuration: to - from,
	}).Result()
	if err != nil {
		return 0, err
	} else if ln := len(avg); ln == 0 {
		return 0, nil
	} else if ln > 1 {
		return 0, fmt.Errorf("%s%d%s", "avg aggregation returned not 1 avg, it is", len(avg), "key="+key)
	}
	return avg[0].Value, nil
}

func (rdb *myRedis) GetCount(ctx context.Context, key string, from, to int) (uint, error) {
	count, err := rdb.TSRangeWithArgs(ctx, key, from, to, &redis.TSRangeOptions{
		Aggregator:     redis.Count,
		BucketDuration: to - from,
	}).Result()
	if err != nil {
		return 0, err
	} else if ln := len(count); ln == 0 {
		return 0, nil
	} else if ln != 1 {
		slog.Error("redis:", "it is not 1 value", count)
		return 0, fmt.Errorf("%s%d", "count aggregation returned not 1 avg, it is ", len(count))
	}
	return uint(count[0].Value), nil
}

func (rdb *myRedis) GetMinimum(ctx context.Context, key string, from, to int) (float64, error) {
	minPrice, err := rdb.TSRangeWithArgs(ctx, key, from, to, &redis.TSRangeOptions{
		Aggregator:     redis.Min,
		BucketDuration: to - from,
	}).Result()
	if err != nil {
		return 0, err
	} else if ln := len(minPrice); ln == 0 {
		return 0, nil
	} else if ln != 1 {
		return 0, fmt.Errorf("%s", "min aggregation returned not 1 avg")
	}
	return minPrice[0].Value, nil
}

func (rdb *myRedis) GetMaximum(ctx context.Context, key string, from, to int) (float64, error) {
	maxPrice, err := rdb.TSRangeWithArgs(ctx, key, from, to, &redis.TSRangeOptions{
		Aggregator:     redis.Max,
		BucketDuration: to - from,
	}).Result()
	if err != nil {
		return 0, err
	} else if ln := len(maxPrice); ln == 0 {
		return 0, nil
	} else if ln != 1 {
		return 0, fmt.Errorf("%s", "max aggregation returned not 1 avg")
	}
	return maxPrice[0].Value, nil
}
