package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"marketflow/internal/domain"
)

func (p *poolDB) GetLowestPriceBySym(ctx context.Context, sym string) (*domain.Exchange, error) {
	const query = `
		SELECT source, max_price, at_time
		FROM exchange_averages
		WHERE symbol = $1
		ORDER BY max_price ASC, at_time DESC
		LIMIT 1;`
	res := &domain.Exchange{
		Symbol: sym,
	}
	var ts time.Time
	err := p.QueryRow(ctx, query, sym).Scan(&res.Source, &res.Price, &ts)
	if err == nil {
		res.Timestamp = ts.UnixMilli()
		return res, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrSymbolNotFound
	}
	return nil, err
}

func (p *poolDB) GetLowestPriceBySymInBackup(ctx context.Context, sym string) (*domain.Exchange, error) {
	const query = `
	SELECT source, price, time_stamp
	FROM exchange_backup
	WHERE symbol = $1
	ORDER BY price ASC, time_stamp DESC
	LIMIT 1;`

	res := &domain.Exchange{
		Symbol: sym,
	}
	var ts time.Time
	err := p.QueryRow(ctx, query, sym).Scan(&res.Source, &res.Price, &ts)
	if err == nil {
		return res, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrSymbolNotFound
	}
	res.Timestamp = ts.UnixMilli()
	return nil, err
}

func (p *poolDB) GetLowestPriceByExSym(ctx context.Context, exName, sym string) (*domain.Exchange, error) {
	const query = `
		SELECT max_price, at_time
		FROM exchange_averages
		WHERE source = $1
			AND symbol = $2
		ORDER BY max_price ASC, at_time DESC
		LIMIT 1;`
	res := &domain.Exchange{
		Source: exName,
		Symbol: sym,
	}
	var ts time.Time
	err := p.QueryRow(ctx, query, exName, sym).Scan(&res.Price, &ts)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrSymbolNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("GetHighestPriceByExSym: %w", err)
	}
	res.Timestamp = ts.UnixMilli()
	return res, nil
}

func (p *poolDB) GetLowestPriceByExSymInBackup(ctx context.Context, ex, sym string) (*domain.Exchange, error) {
	const query = `
	SELECT price, time_stamp
	FROM exchange_backup
	WHERE source=$1 AND symbol = $2
	ORDER BY price ASC, time_stamp DESC
	LIMIT 1;`

	res := &domain.Exchange{
		Source: ex,
		Symbol: sym,
	}
	var ts time.Time
	err := p.QueryRow(ctx, query, sym).Scan(&res.Price, &ts)
	if err == nil {
		return res, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrSymbolNotFound
	}
	res.Timestamp = ts.UnixMilli()
	return nil, err
}
