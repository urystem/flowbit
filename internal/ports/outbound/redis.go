package outbound

import (
	"context"

	"marketflow/internal/domain"
)

type RedisInterForWorkers interface {
	Add(ctx context.Context, ex *domain.Exchange) error
}

type RedisInterGlogal interface {
	// GetAndDelRandomCharacter(ctx context.Context) (*domain.Character, error)
	RedisInterForWorkers
	RedisUseCase
	// GetByLabel(ctx context.Context, from, to int, keys ...string) ([]domain.Exchange, error)
	CloseRedis() error
	RedisForOne
}

type RedisChecker interface {
	CheckHealth(ctx context.Context) error
}

type RedisForOne interface {
	GetAllAverages(ctx context.Context, from, to int) ([]domain.ExchangeAggregation, error)
	RedisChecker
}

type RedisUseCase interface {
	RedisChecker
	latest
	highest
	lowest
	average
}

type latest interface {
	GetLatestPriceBySymbol(ctx context.Context, symbol string) (*domain.Exchange, error)
	GetLastPriceByExAndSym(ctx context.Context, ex, sym string) (*domain.Exchange, error)
}

type highest interface {
	GetHighestPriceWithAlign(ctx context.Context, from int, sym string) (*domain.Exchange, error)
	GetHighestPriceWithEx(ctx context.Context, from int, exName, sym string) (*domain.Exchange, error)
}

type lowest interface {
	GetLowestPriceWithAlign(ctx context.Context, from int, sym string) (*domain.Exchange, error)
	GetLowestPriceWithEx(ctx context.Context, from int, exName, sym string) (*domain.Exchange, error)
}

type average interface {
	GetAveragePriceWithAlign(ctx context.Context, from int, sym string) (*domain.ExchangeAggregation, error)
	GetAveragePriceWithEx(ctx context.Context, from int, exName, sym string) (*domain.ExchangeAggregation, error)
}
