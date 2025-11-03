package postgres

import (
	"context"
	"database/sql"
	"errors"
	"marketflow/internal/domain"
)

func (p *poolDB) GetLatestPriceBySymbol(ctx context.Context, symbol string) (*domain.Exchange, error) {
	const query = `
		SELECT source, price, time_stamp
			FROM exchange_backup
			WHERE symbol = $1
			ORDER BY time_stamp DESC
			LIMIT 1;`
	ex := &domain.Exchange{
		Symbol: symbol,
	}
	err := p.QueryRow(ctx, query, symbol).Scan(&ex.Source, &ex.Price, &ex.Timestamp)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrSymbolNotFound
	}
	return ex, err
}

func (p *poolDB) GetLastPriceByExAndSym(ctx context.Context, ex, sym string) (*domain.Exchange, error) {
	const query = `
		SELECT price, time_stamp
		FROM exchange_backup
		WHERE source = $1 AND symbol = $2
		ORDER BY time_stamp DESC
		LIMIT 1;`
	exchange := &domain.Exchange{
		Source: ex,
		Symbol: sym,
	}
	err := p.QueryRow(ctx, query, sym).Scan(&exchange.Price, &exchange.Timestamp)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrSymbolNotFound
	}
	return exchange, err
}
