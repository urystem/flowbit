package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"marketflow/internal/domain"
	"time"
)

func (p *poolDB) GetHighestPriceBySym(ctx context.Context, sym string) (*domain.Exchange, error) {
	const query = `
		SELECT source, max_price, at_time
		FROM exchange_averages
		WHERE symbol = $1
		ORDER BY max_price DESC, at_time DESC
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

func (p *poolDB) GetHighestPriceBySymInBackup(ctx context.Context, sym string) (*domain.Exchange, error) {
	const query = `
	SELECT source, price, time_stamp
	FROM exchange_backup
	WHERE symbol = $1
	ORDER BY price DESC, time_stamp DESC
	LIMIT 1;`

	res := &domain.Exchange{
		Symbol: sym,
	}
	err := p.QueryRow(ctx, query, sym).Scan(&res.Source, &res.Price, &res.Timestamp)
	if err == nil {
		return res, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrSymbolNotFound
	}
	return nil, err
}

func (p *poolDB) GetHighestPriceBySymWithDuration(ctx context.Context, sym string, from int64) (*domain.Exchange, error) {
	const query = `
		SELECT source, price, time_stamp
		FROM exchange_backup
		WHERE symbol = $1
		  AND time_stamp >= $2
		ORDER BY price DESC, time_stamp DESC
		LIMIT 1;`

	res := &domain.Exchange{Symbol: sym}
	err := p.QueryRow(ctx, query, sym, from).Scan(&res.Source, &res.Price, &res.Timestamp)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrSymbolNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("GetHighestPriceBySymWithDuration: %w", err)
	}
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

// func (p *poolDB) GetHighestPriceBySym(ctx context.Context, sym string) ([]domain.Exchange, error) {
// 	const query = `
// 	WITH max_price_cte AS (
// 		SELECT MAX(max_price) FILTER (WHERE symbol = $1) AS max_price
// 		FROM exchange_averages
// 	)
// 	SELECT DISTINCT ON (ea.source, ea.symbol)
// 			ea.source,
// 			ea.max_price,
// 			ea.at_time
// 	FROM exchange_averages ea
// 	JOIN max_price_cte cte ON ea.max_price = cte.max_price
// 	WHERE ea.symbol = $1
// 	ORDER BY ea.source, ea.at_time DESC;`
// 	rows, err := p.Query(ctx, query, sym)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	var exes []domain.Exchange
// 	for rows.Next() {
// 		var ex domain.Exchange
// 		err := rows.Scan(&ex.Source, &ex.Price, &ex.Timestamp)
// 		if err != nil {
// 			return nil, err
// 		}
// 		exes = append(exes, ex)
// 	}
// 	return exes, nil
// }

// func (p *poolDB) GetHighestPriceBySymInBackup(ctx context.Context, sym string) ([]domain.Exchange, error) {
// 	const query = `
// 	WITH max_price_cte AS (
//   		SELECT MAX(price) FILTER (WHERE symbol = $1) AS max_price
//   		FROM exchange_backup)
// 	SELECT DISTINCT ON (eb.source, eb.symbol)
//   		eb.source,
//   		eb.price,
//   		eb.time_stamp
// 	FROM exchange_backup eb
// 	JOIN max_price_cte cte ON eb.price = cte.max_price
// 	WHERE eb.symbol = $1
// 	ORDER BY eb.source, eb.time_stamp DESC;`

// 	rows, err := p.Query(ctx, query, sym)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	var exes []domain.Exchange
// 	for rows.Next() {
// 		var ex domain.Exchange
// 		err := rows.Scan(&ex.Source, &ex.Price, &ex.Timestamp)
// 		if err != nil {
// 			return nil, err
// 		}
// 		exes = append(exes, ex)
// 	}
// 	return exes, nil
// }
