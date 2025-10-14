package redis

import (
	"context"
	"fmt"
	"strings"

	"marketflow/internal/domain"
)

func (rdb *myRedis) GetAllAverages(ctx context.Context, from, to int) ([]domain.ExchangeAggregation, error) {
	keys, err := rdb.Keys(ctx, "*").Result()
	if err != nil {
		return nil, err
	}

	avgs := make([]domain.ExchangeAggregation, len(keys))
	for i, key := range keys {
		parts := strings.SplitN(key, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("%s%s", "invalid key:", key)
		}
		avgs[i].Source = parts[0]
		avgs[i].Symbol = parts[1]

		avgPrice, err := rdb.GetAveragePrice(ctx, key, from, to)
		if err != nil {
			return nil, err
		}

		avgs[i].AvgPrice = avgPrice

		count, err := rdb.GetCount(ctx, key, from, to)
		if err != nil {
			return nil, err
		}
		avgs[i].Count = count

		minPrice, err := rdb.GetMinimum(ctx, key, from, to)
		if err != nil {
			return nil, err
		}
		avgs[i].MinPrice = minPrice

		maxPrice, err := rdb.GetMaximum(ctx, key, from, to)
		if err != nil {
			return nil, err
		}
		avgs[i].MaxPrice = maxPrice
	}
	return avgs, nil
}
