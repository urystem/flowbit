package redis

import (
	"context"
	"fmt"

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
	} else if len(avg) != 1 {
		fmt.Println(avg, to-from)
		return 0, fmt.Errorf("%s%d", "avg aggregation returned not 1 avg, it is", len(avg))
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
	} else if len(count) != 1 {
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
	} else if len(minPrice) != 1 {
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
	} else if len(maxPrice) != 1 {
		return 0, fmt.Errorf("%s", "max aggregation returned not 1 avg")
	}
	return maxPrice[0].Value, nil
}
