package outbound

import (
	"marketflow/internal/domain"
)

type RedisInterLocal interface {
	Start(exCh <-chan *domain.Exchange, fallbackCh chan<- *domain.Exchange)
}

type RedisInterGlogal interface {
	// GetAndDelRandomCharacter(ctx context.Context) (*domain.Character, error)
	RedisInterLocal
	CloseRedis() error
}
