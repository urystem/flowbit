package outbound

import (
	"context"

	"marketflow/internal/domain"
)

type RedisInterForWorkers interface {
	Add(ctx context.Context, ex *domain.Exchange) error
	GetByLabel(ctx context.Context, from, to int, keys ...string) ([]domain.Exchange, error)
	GetAvarages(ctx context.Context) ([]domain.ExchangeAvg, error)
}

type RedisInterGlogal interface {
	// GetAndDelRandomCharacter(ctx context.Context) (*domain.Character, error)
	RedisInterForWorkers
	CloseRedis() error
}
