package app

// import (
// 	"context"
// 	"fmt"
// 	"log/slog"
// 	"time"

// 	"marketflow/internal/services/batcher"
// )

// // fall керек емес сияктыгой
// func (app *myApp) timerOneMinute(fall batcher.InsertAndStatus) {
// 	const interval = time.Minute
// 	from := time.Now().Truncate(interval)
// 	next := from.Add(interval) // ближайшая "ровная" минута
// 	timer := time.NewTimer(time.Until(next))
// 	defer slog.Info("timer stoped")
// 	defer timer.Stop()
// 	for {
// 		select {
// 		case <-app.ctx.Done():
// 			return
// 		case <-timer.C:
// 			ctx, cancel := context.WithTimeout(app.ctx, 27*time.Second)
// 			avgs, err := app.red.GetAllAverages(ctx, int(from.UnixMilli()), int(next.UnixMilli()))
// 			cancel()
// 			if err != nil {
// 				slog.Error("ticker", "error:", err)
// 			} else {
// 				ctx, cancel := context.WithTimeout(app.ctx, 27*time.Second)
// 				err = app.db.SaveWithCopyFrom(ctx, avgs, from)
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

// func (app *myApp) tickerToCheckFallBack() {
// 	batch := make([]domain.Exchange, 0, 512)

// 	for ex := range app.fallBack {

// 		batch = append(batch, *ex)
// 		if len(batch) == 512 {
// 			fmt.Println("batched")
// 			batch = batch[:0]
// 		}
// 	}
// }
