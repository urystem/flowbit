package batcher

import (
	"log/slog"

	"marketflow/internal/domain"
)

func (f *batchCollector) GoAndReturnCh() chan<- *domain.Exchange {
	go f.goFunc()
	return f.channel
}

func (f *batchCollector) IsNotWorking() bool {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	return f.sendedSignalNotWorking
}

func (f *batchCollector) InsertBatches() error {
	err := f.sql.FallBack(f.ctx, f.batch)
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
