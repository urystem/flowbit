package redis

import (
	"context"
	"fmt"
	"marketflow/internal/domain"
)

func (rdb *myRedis) GetLatestPriceBySymbol(ctx context.Context, symbol string) (*domain.Exchange, error) {
	res, err := rdb.TSMGet(ctx, []string{"symbol=" + symbol}).Result()
	if err != nil {
		return nil, err
	} else if len(res) == 0 {
		return nil, domain.ErrSymbolNotFound
	}
	ex := &domain.Exchange{
		Symbol: symbol,
	}

	for exSym, v := range res {
		if len(v) != 2 {
			return nil, fmt.Errorf("%s", "unknown error (getlatest)")
		}
		timePriceSlc, ok := v[1].([]any)
		if !ok {
			return nil, fmt.Errorf("%s", "unkown error (not []any)")
		}

		currTime, ok := timePriceSlc[0].(int64)
		if !ok {
			return nil, fmt.Errorf("%s", "invalid time")
		}

		thisPrice, ok := timePriceSlc[1].(float64)
		if !ok {
			return nil, fmt.Errorf("%s", "invalid price")
		}
		if ex.Timestamp < currTime {
			ex.Source = exSym
			ex.Timestamp = currTime
			ex.Price = thisPrice
		}
	}
	return ex, nil
}

func (rdb *myRedis) GetLastPriceByExAndSym(ctx context.Context, ex, sym string) (*domain.Exchange, error) {
	res, err := rdb.TSGet(ctx, ex+":"+sym).Result()
	if err != nil {
		return nil, err
	}
	if res.Timestamp == 0 {
		return nil, domain.ErrSymbolNotFound
	}
	return &domain.Exchange{
		Source:    ex,
		Symbol:    sym,
		Price:     res.Value,
		Timestamp: res.Timestamp,
	}, nil
}
