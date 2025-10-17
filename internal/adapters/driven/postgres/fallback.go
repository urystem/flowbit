package postgres

import (
	"context"
	"time"

	"marketflow/internal/domain"

	"github.com/jackc/pgx/v5"
)

func (db *poolDB) FallBack(ctx context.Context, exs []*domain.Exchange) error {
	rows := make([][]any, len(exs))
	for i, v := range exs {
		rows[i] = []any{
			v.Source,
			v.Symbol,
			v.Price,
			time.UnixMilli(v.Timestamp),
		}
	}
	_, err := db.CopyFrom(
		ctx,
		pgx.Identifier{"exchange_backup"},
		[]string{"source", "symbol", "price", "time_stamp"},
		pgx.CopyFromRows(rows),
	)
	return err
}
