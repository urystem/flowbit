package inbound

import (
	"context"
	"marketflow/internal/domain"
)

type UsecaseInter interface {
	SwitchToTest()
	SwitchToLive()
	CheckHealth(ctx context.Context) any
	GetLatestBySymbol(ctx context.Context, symbol string) (*domain.Exchange, error)
	GetLatestPriceByExAndSym(context.Context, string, string) (*domain.Exchange, error)
	GetHighestPriceBySym(ctx context.Context, sym string) (*domain.Exchange, error)
}
