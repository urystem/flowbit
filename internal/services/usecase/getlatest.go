package usecase

import (
	"context"
)

// if redis working, we can get latest in redis

/*
if redis not working, we get latest in sql
but we cannot control batch slice, so may be we will get NULL(no any exchanges).
Because batches not fully inserted to sql
*/
func (u *myUsecase) GetLatestBySymbol(ctx context.Context, symbol string) (float64, error) {
	if !u.one.IsNotWorking() {
		return u.rdb.GetLatestPriceBySymbol(ctx, symbol)
	}
	return u.db.GetLatestPriceBySymbol(ctx, symbol)
}

func (u *myUsecase) GetLatestPriceByExAndSym(ctx context.Context, ex, sym string) (float64, error) {
	if !u.one.IsNotWorking() {
		return u.rdb.GetLastPriceByExAndSym(ctx, ex, sym)
	}
	return u.db.GetLastPriceByExAndSym(ctx, ex, sym)
}
