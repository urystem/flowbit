package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"marketflow/internal/domain"
	"time"
)

func (p *poolDB) GetAveragePriceBySymTime(ctx context.Context, sym string, from time.Time) (*domain.ExchangeAggregation, error) {
	const query = `
	SELECT
  		SUM(count) AS total_count,
  		SUM(average_price * count) / SUM(count) AS weighted_avg,
  		MIN(at_time) AS first_time
	FROM exchange_averages
	WHERE symbol = $1
	AND at_time >= $2
	HAVING SUM(count) > 0;`

	var ex domain.ExchangeAggregation
	if err := p.QueryRow(ctx, query, sym, from).Scan(&ex.Count, &ex.AvgPrice, &ex.Timestamp); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrSymbolNotFound
		}
		return nil, fmt.Errorf("GetAveragePriceBySym: %w", err)
	}
	return &ex, nil
}

func (p *poolDB) GetAveragePriceBySymInBackupTime(ctx context.Context, sym string, from time.Time) (*domain.ExchangeAggregation, error) {
	const query = `
	SELECT 
		COUNT(*) AS total_count,
		AVG(price) AS average_price,
  		MIN(time_stamp) AS first_time
	FROM exchange_backup
	WHERE symbol = $1
		AND time_stamp >= $2
	HAVING COUNT(*) > 0;`
	var ex domain.ExchangeAggregation
	if err := p.QueryRow(ctx, query, sym, from).Scan(&ex.Count, &ex.AvgPrice, &ex.Timestamp); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrSymbolNotFound
		}
		return nil, fmt.Errorf("GetAveragePriceBySym: %w", err)
	}
	return &ex, nil
}

func (p *poolDB) GetAveragePriceByExSymTime(ctx context.Context, exName, sym string, from time.Time) (*domain.ExchangeAggregation, error) {
	const query = `
	SELECT
  		SUM(count) AS total_count,
  		SUM(average_price * count) / SUM(count) AS weighted_avg,
		MIN(at_time) AS first_time
	FROM exchange_averages
	WHERE  source = $1
	AND symbol = $2
	AND at_time >= $3
	HAVING SUM(count) > 0;`

	ex := &domain.ExchangeAggregation{
		Source: exName,
		Symbol: sym,
	}
	if err := p.QueryRow(ctx, query, exName, sym, from).Scan(&ex.Count, &ex.AvgPrice, &ex.Timestamp); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrSymbolNotFound
		}
		return nil, fmt.Errorf("GetAveragePriceBySym: %w", err)
	}
	return ex, nil
}

func (p *poolDB) GetAveragePriceByExSymInBackupTime(ctx context.Context, exName, sym string, from time.Time) (*domain.ExchangeAggregation, error) {
	const query = `
	SELECT 
		COUNT(*) AS total_count,
		AVG(price) AS average_price,
  		MIN(time_stamp) AS first_time
	FROM exchange_backup
	WHERE source = $1
		AND symbol = $2
		AND time_stamp >= $3
	HAVING COUNT(*) > 0;`
	ex := &domain.ExchangeAggregation{
		Source: exName,
		Symbol: sym,
	}
	if err := p.QueryRow(ctx, query, exName, sym, from).Scan(&ex.Count, &ex.AvgPrice, &ex.Timestamp); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrSymbolNotFound
		}
		return nil, fmt.Errorf("GetAveragePriceBySym: %w", err)
	}
	return ex, nil
}
