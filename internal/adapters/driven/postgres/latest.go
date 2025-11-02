package postgres

import (
	"context"
	"database/sql"
	"errors"
	"marketflow/internal/domain"
)

func (p *poolDB) GetLatestPriceBySymbol(ctx context.Context, symbol string) (float64, error) {
	const query = `
		SELECT price
			FROM exchange_backup
			WHERE symbol = $1
			ORDER BY time_stamp DESC
			LIMIT 1;`
	var price float64
	err := p.QueryRow(ctx, query, symbol).Scan(&price)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, domain.ErrSymbolNotFound
	}
	return price, err
}

func (p *poolDB) GetLastPriceByExAndSym(ctx context.Context, ex, sym string) (float64, error) {
	const query = `
		SELECT price
		FROM exchange_backup
		WHERE source = $1 AND symbol = $2
		ORDER BY time_stamp DESC
		LIMIT 1;`
	var price float64
	err := p.QueryRow(ctx, query, sym).Scan(&price)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, domain.ErrSymbolNotFound
	}
	return price, err
}
