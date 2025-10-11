package postgres

import (
	"context"

	"marketflow/internal/domain"

	"github.com/jackc/pgx/v5"
)

// import "context"

func (db *poolDB) SaveAverage(ctx context.Context, avgs []domain.ExchangeAvg) error {
	const sql = `
	INSERT INTO exchange_averages
	(source, symbol, count, average_price, min_price, max_price, at_time)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // откат, если что-то пойдёт не так

	batch := &pgx.Batch{}
	for _, v := range avgs {
		batch.Queue(sql,
			v.Source,
			v.Symbol,
			v.Count,
			v.AvgPrice,
			v.MinPrice,
			v.MaxPrice,
			v.AtTime,
		)
	}
	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	for range avgs {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (db *poolDB) SaveWithCopyFrom(ctx context.Context, avgs []domain.ExchangeAvg) error {
	rows := make([][]any, len(avgs))
	for i, v := range avgs {
		rows[i] = []any{
			v.Source,
			v.Symbol,
			v.Count,
			v.AvgPrice,
			v.MinPrice,
			v.MaxPrice,
			v.AtTime,
		}
	}

	_, err := db.CopyFrom(
		ctx,
		pgx.Identifier{"exchange_averages"},
		[]string{"source", "symbol", "count", "average_price", "min_price", "max_price", "at_time"},
		pgx.CopyFromRows(rows),
	)
	return err
}
