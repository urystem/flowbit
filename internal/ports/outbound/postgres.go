package outbound

import (
	"context"
	"time"

	"marketflow/internal/domain"
)

type PgxInter interface {
	CloseDB()
	PgxForTimerAndBatcher
	PgxForUseCase
	DBTestCleaner
}

type PgxForTimerAndBatcher interface {
	GetAverageAndDelete(ctx context.Context, from, to time.Time) ([]domain.ExchangeAggregation, error)
	SaveWithCopyFrom(ctx context.Context, avgs []domain.ExchangeAggregation, ti time.Time) error
	FallBack(ctx context.Context, rows [][]any) error
	CleanerTest(ctx context.Context) error
}

type PgxForUseCase interface {
	CheckHealth(ctx context.Context) error
	pgxLatest
	pgxHighest
	pgxLowest
	pgxAverage
}

type pgxLatest interface {
	GetLatestPriceBySymbol(ctx context.Context, symbol string) (*domain.Exchange, error)
	GetLastPriceByExAndSym(ctx context.Context, ex, sym string) (*domain.Exchange, error)
}

type pgxHighest interface {
	GetHighestPriceBySym(ctx context.Context, sym string) (*domain.Exchange, error)
	GetHighestPriceBySymInBackup(ctx context.Context, sym string) (*domain.Exchange, error)
	GetHighestPriceBySymWithDuration(ctx context.Context, sym string, from time.Time) (*domain.Exchange, error)
	GetHighestPriceBySymWithDurationInAverage(ctx context.Context, sym string, from time.Time) (*domain.Exchange, error)

	GetHighestPriceByExSym(ctx context.Context, exName, sym string) (*domain.Exchange, error)
	GetHighestPriceByExSymInBackup(ctx context.Context, ex, sym string) (*domain.Exchange, error)
	GetHighestPriceByExSymWithDurationInAverage(ctx context.Context, ex, sym string, from time.Time) (*domain.Exchange, error)
	GetHighestPriceByExSymWithDuration(ctx context.Context, ex, sym string, from time.Time) (*domain.Exchange, error)
}

type pgxLowest interface {
	GetLowestPriceBySym(ctx context.Context, sym string) (*domain.Exchange, error)
	GetLowestPriceBySymInBackup(ctx context.Context, sym string) (*domain.Exchange, error)
	GetLowestPriceByExSym(ctx context.Context, exName, sym string) (*domain.Exchange, error)
	GetLowestPriceByExSymInBackup(ctx context.Context, ex, sym string) (*domain.Exchange, error)

	GetLowestPriceBySymWithDuration(ctx context.Context, sym string, from time.Time) (*domain.Exchange, error)
	GetLowestPriceBySymWithDurationInAverage(ctx context.Context, sym string, from time.Time) (*domain.Exchange, error)
	GetLowestPriceByExSymWithDurationInAverage(ctx context.Context, ex, sym string, from time.Time) (*domain.Exchange, error)
	GetLowestPriceByExSymWithDuration(ctx context.Context, ex, sym string, from time.Time) (*domain.Exchange, error)
}

type pgxAverage interface {
	GetAveragePriceBySym(ctx context.Context, sym string) (*domain.ExchangeAggregation, error)
	GetAveragePriceBySymInBackup(ctx context.Context, sym string) (*domain.ExchangeAggregation, error)
	GetAveragePriceByExSym(ctx context.Context, exName, sym string) (*domain.ExchangeAggregation, error)
	GetAveragePriceByExSymInBackup(ctx context.Context, exName, sym string) (*domain.ExchangeAggregation, error)

	GetAveragePriceBySymTime(ctx context.Context, sym string, from time.Time) (*domain.ExchangeAggregation, error)
	GetAveragePriceBySymInBackupTime(ctx context.Context, sym string, from time.Time) (*domain.ExchangeAggregation, error)
	GetAveragePriceByExSymTime(ctx context.Context, exName, sym string, from time.Time) (*domain.ExchangeAggregation, error)
	GetAveragePriceByExSymInBackupTime(ctx context.Context, exName, sym string, from time.Time) (*domain.ExchangeAggregation, error)
}

type DBTestCleaner interface {
	CleanerTest(ctx context.Context) error
}
