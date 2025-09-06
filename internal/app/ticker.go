package app

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
