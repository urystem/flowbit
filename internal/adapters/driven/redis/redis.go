package redis

import (
	"context"
	"log/slog"

	"marketflow/internal/domain"
	"marketflow/internal/ports/inbound"
	"marketflow/internal/ports/outbound"

	"github.com/redis/go-redis/v9"
)

type myRedis struct {
	// log
	// errCh
	ctx context.Context
	*redis.Client
}

func InitRickRedis(ctx context.Context, red inbound.RedisConfig) (outbound.RedisInterGlogal, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redisDB:" + red.GetAddr(), // имя сервиса + порт                 // адрес Redis
		Password: red.GetPass(),              // пароль, если есть
		// DB:       0,                          // номер БД (0 по умолчанию)
	})
	return &myRedis{ctx, rdb}, rdb.Ping(ctx).Err()
}

func (rdb *myRedis) Start(exCh <-chan *domain.Exchange, fallbackCh chan<- *domain.Exchange) {
	go func() {
		defer slog.Info("redis stopped")

		for {
			select {
			case <-rdb.ctx.Done():
				// Контекст отменён — выходим
				close(fallbackCh)
				return

			case ex, ok := <-exCh:
				if !ok {
					// Канал закрыт — больше данных не будет
					close(fallbackCh)
					return
				}
				_, err := rdb.TSAddWithArgs(
					rdb.ctx,
					ex.Symbol,
					ex.Timestamp,
					ex.Price,
					&redis.TSOptions{},
				).Result()
				if err != nil {
					// Если Redis упал — fallback в Postgres
					fallbackCh <- ex
				}
			}
		}
	}()
}

func (rdb *myRedis) CloseRedis() error {
	return rdb.Close()
}
