package app

import (
	"log/slog"
	"time"
)

func (app *myApp) tickerOneMinute() {
	const interval = time.Minute
	next := time.Now().Truncate(interval).Add(interval) // ближайшая "ровная" минута
	timer := time.NewTimer(time.Until(next))
	defer timer.Stop()
	for {
		select {
		case <-app.ctx.Done():
			return
		case <-timer.C:
			avgs, err := app.red.GetAvarages2(app.ctx, int(next.UnixMilli()))
			if err != nil {
				slog.Error("ticker", "error:", err)
			} else {
				err = app.db.SaveWithCopyFrom(app.ctx, avgs)
				if err != nil {
					slog.Error("ticker", "psql", err)
				} else {
					slog.Info("saved to sql")
				}
			}

			next = next.Add(interval)
			timer.Reset(time.Until(next))
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
