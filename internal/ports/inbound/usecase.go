package inbound

import (
	"context"
	"time"

	"marketflow/internal/domain"
)

type UsecaseInter interface {
	CheckHealth(ctx context.Context) any
	mode
	latest
	highest
	lowest
	average
}

type mode interface {
	SwitchToTest()
	SwitchToLive()
}

type highest interface {
	GetHighestPriceBySym(ctx context.Context, sym string) (*domain.Exchange, error)
	GetHighestPriceBySymWithDuration(ctx context.Context, sym string, duration time.Duration) (any, error)
	GetHighestPriceByExSym(ctx context.Context, exName, sym string) (*domain.Exchange, error)
	GetHighestPriceByExSymDuration(ctx context.Context, exName, sym string, dur time.Duration) (any, error)
}

type latest interface {
	GetLatestBySymbol(ctx context.Context, symbol string) (*domain.Exchange, error)
	GetLatestPriceByExAndSym(ctx context.Context, ex, sym string) (*domain.Exchange, error)
}

type lowest interface {
	GetLowestPriceBySym(ctx context.Context, sym string) (*domain.Exchange, error)
	GetLowestPriceBySymWithDuration(ctx context.Context, sym string, duration time.Duration) (any, error)
	GetLowestPriceByExSym(ctx context.Context, exName, sym string) (*domain.Exchange, error)
	GetLowestPriceByExSymDuration(ctx context.Context, exName, sym string, dur time.Duration) (any, error)
}

type average interface {
	GetAveragePriceBySym(ctx context.Context, sym string) (*domain.ExchangeAggregation, error)
	GetAveragePriceBySymWithDuration(ctx context.Context, sym string, dur time.Duration) (*domain.ExchangeAggregation, error)
	GetAveragePriceByExSym(ctx context.Context, exName, sym string) (*domain.ExchangeAggregation, error)
	GetAveragePriceByExSymDuration(ctx context.Context, exName, sym string, dur time.Duration) (any, error)
}
