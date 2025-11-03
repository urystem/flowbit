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
	avgs := make([]domain.ExchangeAggregation, 0, len(keys))
	for _, key := range keys {
		parts := strings.SplitN(key, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("%s%s", "invalid key:", key)
		}

		count, err := rdb.GetCount(ctx, key, from, to)
		if err != nil {
			return nil, err
		} else if count == 0 {
			continue
		}

		avgPrice, err := rdb.GetAveragePrice(ctx, key, from, to)
		if err != nil {
			return nil, err
		}

		minPrice, err := rdb.GetMinimum(ctx, key, from, to)
		if err != nil {
			return nil, err
		}

		maxPrice, err := rdb.GetMaximum(ctx, key, from, to)
		if err != nil {
			return nil, err
		}

		avgs = append(avgs, domain.ExchangeAggregation{
			Source:   parts[0],
			Symbol:   parts[1],
			Count:    count,
			AvgPrice: avgPrice,
			MinPrice: minPrice,
			MaxPrice: maxPrice,
		})
	}
	return avgs, nil
}
