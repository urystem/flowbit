package postgres

import (
	"context"
	"fmt"

	"marketflow/internal/ports/inbound"
	"marketflow/internal/ports/outbound"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB
type poolDB struct {
	*pgxpool.Pool
}

func InitDB(ctx context.Context, cfg inbound.DBConfig) (outbound.PostGres, error) {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.GetUser(),
		cfg.GetPassword(),
		cfg.GetHostName(),
		cfg.GetPort(),
		cfg.GetDBName(),
		cfg.GetSSLMode(),
	)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &poolDB{pool}, nil
}

func (pool *poolDB) CloseDB() {
	pool.Close()
}

