package outbound

import (
	"context"

	"marketflow/internal/domain"
)

type RedisInterLocal interface {
	SetExchange(ctx context.Context, ex *domain.Exchange) error
}

type RedisInterGlogal interface {
	// GetAndDelRandomCharacter(ctx context.Context) (*domain.Character, error)
	RedisInterLocal
	CloseRedis() error
}
