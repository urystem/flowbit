package app

import (
	"fmt"
	"log/slog"
	"time"

	"marketflow/internal/domain"
)

func (app *myApp) timerOneMinute() {
	const interval = time.Minute
	from := time.Now().Truncate(interval)
	next := from.Add(interval) // ближайшая "ровная" минута
	timer := time.NewTimer(time.Until(next))
	defer timer.Stop()
	for {
		select {
		case <-app.ctx.Done():
			return
		case <-timer.C:

			avgs, err := app.red.GetAllAverages(app.ctx, int(from.UnixMilli()), int(next.UnixMilli()))
			if err != nil {
				slog.Error("ticker", "error:", err)
			} else {
				err = app.db.SaveWithCopyFrom(app.ctx, avgs, from)
				if err != nil {
					slog.Error("ticker", "psql", err)
				} else {
					slog.Info("saved to sql")
				}
			}
			from = next
			next = next.Add(interval)
			timer.Reset(time.Until(next))
		}
	}
}

func (app *myApp) tickerToCheckFallBack() {
	batch := make([]domain.Exchange, 0, 512)

	for ex := range app.fallBack {

		batch = append(batch, *ex)
		if len(batch) == 512 {
			fmt.Println("batched")
			batch = batch[:0]
		}
	}
}

// func (app *myApp) initTicker() {
// 	app.wg.Add(2)

// 	signal := make(chan struct{})
// 	ticker := time.NewTicker(1 * time.Minute) // 🕒 каждый 1 минуту

// 	app.srv.RegisterOnShutDown(func() {
// 		defer app.wg.Done()
// 		ticker.Stop()
// 		signal <- struct{}{}
// 	})

// 	go func() {
// 		defer app.wg.Done()
// 		for {
// 			select {
// 			case <-signal:
// 				return
// 			case <-ticker.C:
// 				app.tickerToDo()
// 			}
// 		}
// 	}()
// }

// func (app *myApp) tickerToDo() {
// 	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
// 	defer cancel()
// 	err := app.ticker.Archiver(ctx)
// 	if err != nil {
// 		slog.Error(err.Error())
// 	}
// }
