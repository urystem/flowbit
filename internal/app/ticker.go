package app

import (
	"log/slog"
	"time"
)

func (app *myApp) tickerOneMinute() {
	ticker := time.NewTicker(1 * time.Minute) // 🕒 каждый 1 минуту
	defer ticker.Stop()
	for {
		select {
		case <-app.ctx.Done():
			return
		case <-ticker.C:
			avgs, err := app.red.GetAvarages(app.ctx)
			if err != nil {
				slog.Error("ticker", "error:", err)
				continue
			}
			err = app.db.SaveWithCopyFrom(app.ctx, avgs)
			if err != nil {
				slog.Error("ticker", "psql", err)
			} else {
				slog.Info("saved to sql")
			}
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
