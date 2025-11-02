package postgres

import (
	"context"
	"database/sql"
	"errors"
	"marketflow/internal/domain"
)

func (p *poolDB) GetHighestPriceBySym(ctx context.Context, sym string) (float64, error) {
	const query = `
		SELECT MAX(max_price)
		FROM exchange_averages
		WHERE symbol = $1;`

	var price float64
	err := p.QueryRow(ctx, query, sym).Scan(&price)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, domain.ErrSymbolNotFound
	}
	return price, err
}

func (p *poolDB) GetHighestPriceBySymInBackup(ctx context.Context, sym string) (float64, error) {
	const query = `
		SELECT MAX(price)
		FROM exchange_backup
		WHERE symbol = $1;`

	var price float64
	err := p.QueryRow(ctx, query, sym).Scan(&price)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, domain.ErrSymbolNotFound
	}
	return price, err
}
