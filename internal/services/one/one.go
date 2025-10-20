package one

import (
	"marketflow/internal/ports/outbound"
	"marketflow/internal/services/batcher"
)

type oneMinute struct {
	wasErr bool // there was error in this last minute
	fall   batcher.InsertAndStatus
	red    outbound.RedisForOne
	db     outbound.PgxForTimer
}

func NewTimerOneMinute(fall batcher.InsertAndStatus, red outbound.RedisForOne, db outbound.PgxForTimer) {
}

// 	const interval = time.Minute
// 	from := time.Now().Truncate(interval)
// 	next := from.Add(interval) // ближайшая "ровная" минута
// 	timer := time.NewTimer(time.Until(next))
// 	defer slog.Info("timer stoped")
// 	defer timer.Stop()
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		case <-timer.C:
// 			ctx, cancel := context.WithTimeout(ctx, 27*time.Second)
// 			avgs, err := red.GetAllAverages(ctx, int(from.UnixMilli()), int(next.UnixMilli()))
// 			cancel()
// 			if err != nil {
// 				slog.Error("ticker", "error:", err)
// 			} else {
// 				ctx, cancel := context.WithTimeout(ctx, 27*time.Second)
// 				err = db.SaveWithCopyFrom(ctx, avgs, from)
// 				cancel()
// 				if err != nil {
// 					slog.Error("ticker", "psql", err)
// 				} else {
// 					fmt.Println(avgs)
// 					slog.Info("saved to sql")
// 				}
// 			}
// 			from = next
// 			next = next.Add(interval)
// 			timer.Reset(time.Until(next))
// 		}
// 	}
// }

// skipper
