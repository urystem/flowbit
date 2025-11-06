package redis

import (
	"context"
	"fmt"
)

func (r *myRedis) TestCleaner(ctx context.Context) error {
	keys, err := r.Keys(ctx, "test*").Result()
	if err != nil {
		return fmt.Errorf("ошибка получения ключей: %w", err)
	}

	// Удаляем все найденные ключи сразу
	if err := r.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("ошибка удаления: %w", err)
	}
	return nil
}
