package batcher

import (
	"context"
	"log/slog"

	"marketflow/internal/domain"
)

func (f *batchCollector) GoAndReturnCh() chan<- *domain.Exchange {
	go f.goFunc()
	return f.channel
}

func (f *batchCollector) IsNotWorking() bool {
	return f.notWorking.Load()
}

func (f *batchCollector) InsertBatches(ctx context.Context) error {
	err := f.sql.FallBack(ctx, f.batch)
	if err != nil {
		slog.Error("fallback", "sql потеря данных", err)
	} else {
		slog.Info("batched")
	}
	for i := range f.batch {
		f.put(f.batch[i])
	}
	f.batch = f.batch[:0]
	return err
}

