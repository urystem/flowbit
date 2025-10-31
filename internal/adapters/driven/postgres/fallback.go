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

func (db *poolDB) GetAverageAndDelete(ctx context.Context, from, to time.Time) ([]domain.ExchangeAggregation, error) {
	const query = `
	WITH deleted AS (
		DELETE FROM exchange_backup
		WHERE time_stamp BETWEEN $1 AND $2
		RETURNING source, symbol, price
	)
	SELECT
		source,
		symbol,
		COUNT(*) AS count,
		AVG(price) AS avg_price,
		MIN(price) AS min_price,
		MAX(price) AS max_price
	FROM deleted
	GROUP BY source, symbol;`

	rows, err := db.Query(ctx, query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.ExchangeAggregation
	for rows.Next() {
		var agg domain.ExchangeAggregation
		var count int64
		if err := rows.Scan(&agg.Source, &agg.Symbol, &count, &agg.AvgPrice, &agg.MinPrice, &agg.MaxPrice); err != nil {
			return nil, err
		}
		agg.Count = uint(count)
		result = append(result, agg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
