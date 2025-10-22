package one

import (
	"time"
)

func (one *oneMinute) goFuncTimer() {
	from := time.Now().Truncate(time.Minute)
	next := from.Add(time.Minute) // ближайшая "ровная" минута
	timer := time.NewTimer(time.Until(next))
	// defer slog.Info("timer stoped")
	defer timer.Stop()
	for {
		select {
		case <-one.ctx.Done():
			return
		case <-timer.C:
			one.insertAverage(one.ctx, from, next)
			from = next
			next = next.Add(time.Minute)
			timer.Reset(time.Until(next))
		}
	}
}
