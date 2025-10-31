package inbound

import "context"

type UsecaseInter interface {
	SwitchToTest()
	SwitchToLive()
	CheckHealth(ctx context.Context) any
	GetLatestBySymbol(ctx context.Context, symbol string) (float64, error)
	GetLatestPriceByExAndSym(context.Context, string, string) (float64, error)
}
