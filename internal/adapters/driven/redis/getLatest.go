package redis

import (
	"context"
	"fmt"
)

func (rdb *myRedis) GetLatestPriceBySymbol(ctx context.Context, symbol string) (float64, error) {
	res, err := rdb.TSMGet(ctx, []string{"symbol=" + symbol}).Result()
	if err != nil {
		return 0, err
	} else if len(res) == 0 {
		return 0, fmt.Errorf("symbol not found")
	}
	var ans float64
	var lastTime int64
	for _, v := range res {
		if len(v) != 2 {
			return 0, fmt.Errorf("%s", "unknown error (getlatest)")
		}
		timePriceSlc, ok := v[1].([]any)
		if !ok {
			return 0, fmt.Errorf("%s", "unkown error (not []any)")
		}

		currTime, ok := timePriceSlc[0].(int64)
		if !ok {
			return 0, fmt.Errorf("%s", "invalid time")
		}

		thisPrice, ok := timePriceSlc[1].(float64)
		if !ok {
			return 0, fmt.Errorf("%s", "invalid price")
		}
		if lastTime < currTime {
			lastTime = currTime
			ans = thisPrice
		}
	}
	return ans, nil
}

func (rdb *myRedis) GetLastPriceByExAndSym(ctx context.Context, ex, sym string) (float64, error) {
	res, err := rdb.TSGet(ctx, ex+":"+sym).Result()
	if err != nil {
		return 0, err
	}
	return res.Value, nil
}
