package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"marketflow/internal/domain"
)

func (p *poolDB) GetHighestPriceBySymWithDuration(ctx context.Context, sym string, from time.Time) (*domain.Exchange, error) {
	const query = `
		SELECT source, price, time_stamp
		FROM exchange_backup
		WHERE symbol = $1
		  AND time_stamp >= $2
		ORDER BY price DESC, time_stamp DESC
		LIMIT 1;`

	res := &domain.Exchange{Symbol: sym}
	var ts time.Time
	err := p.QueryRow(ctx, query, sym, from).Scan(&res.Source, &res.Price, &ts)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrSymbolNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("GetHighestPriceBySymWithDuration: %w", err)
	}
	res.Timestamp = ts.UnixMilli()
	return res, nil
}

func (p *poolDB) GetHighestPriceBySymWithDurationInAverage(ctx context.Context, sym string, from time.Time) (*domain.Exchange, error) {
	const query = `
		SELECT source, max_price, at_time
		FROM exchange_averages
		WHERE symbol = $1
		  AND at_time >= $2
		ORDER BY max_price DESC, at_time DESC
		LIMIT 1;`

	res := &domain.Exchange{Symbol: sym}
	var ts time.Time
	err := p.QueryRow(ctx, query, sym, from).Scan(&res.Source, &res.Price, &ts)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrSymbolNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("GetHighestPriceBySymWithDurationInAverage: %w", err)
	}
	res.Timestamp = ts.UnixMilli()
	return res, nil
}

func (p *poolDB) GetHighestPriceByExSymWithDurationInAverage(ctx context.Context, ex, sym string, from time.Time) (*domain.Exchange, error) {
	const query = `
		SELECT max_price, at_time
		FROM exchange_averages
		WHERE source = $1
			AND symbol = $2
		  	AND at_time >= $3
		ORDER BY max_price DESC, at_time DESC
		LIMIT 1;`

	res := &domain.Exchange{
		Source: ex,
		Symbol: sym,
	}
	var ts time.Time
	err := p.QueryRow(ctx, query, ex, sym, from).Scan(&res.Price, &ts)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrSymbolNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("GetHighestPriceBySymWithDurationInAverage: %w", err)
	}
	res.Timestamp = ts.UnixMilli()
	return res, nil
}

func (p *poolDB) GetHighestPriceByExSymWithDuration(ctx context.Context, ex, sym string, from time.Time) (*domain.Exchange, error) {
	const query = `
		SELECT price, time_stamp
		FROM exchange_backup
		WHERE source=$1 
			AND symbol = $2
		  	AND time_stamp >= $3
		ORDER BY price DESC, time_stamp DESC
		LIMIT 1;`

	res := &domain.Exchange{
		Source: ex,
		Symbol: sym,
	}
	var ts time.Time
	err := p.QueryRow(ctx, query, ex, sym, from).Scan(&res.Price, &ts)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrSymbolNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("GetHighestPriceBySymWithDuration: %w", err)
	}
	res.Timestamp = ts.UnixMilli()
	return res, nil
}
