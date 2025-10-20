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
	RedisChecker
	RedisInterForWorkers
	GetByLabel(ctx context.Context, from, to int, keys ...string) ([]domain.Exchange, error)
	CloseRedis() error
	RedisForOne
}

type RedisChecker interface {
	CheckHealth(ctx context.Context) error
}

type RedisForOne interface {
	GetAllAverages(ctx context.Context, from, to int) ([]domain.ExchangeAggregation, error)
}
