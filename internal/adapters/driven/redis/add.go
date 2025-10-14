package redis

import (
	"context"

	"marketflow/internal/domain"

	"github.com/redis/go-redis/v9"
)

func (rdb *myRedis) Add(ctx context.Context, ex *domain.Exchange) error {
	_, err := rdb.TSAddWithArgs(
		ctx,
		ex.Source+":"+ex.Symbol, // tenge
		ex.Timestamp,            // time
		ex.Price,                // price
		&redis.TSOptions{
			Retention:       70000, // 70 SEC
			DuplicatePolicy: "LAST",
			Labels: map[string]string{
				"exchange": ex.Source,
				"symbol":   ex.Symbol,
			},
		},
	).Result()
	return err
}
