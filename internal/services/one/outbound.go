package one

import (
	"context"
	"fmt"
)

func (one *oneMinute) IsNotWorking() bool {
	return one.notWorking.Load()
}

func (one *oneMinute) Run(ctx context.Context) error {
	if one.ctx != nil {
		return fmt.Errorf("%s", "already")
	}
	one.ctx = ctx
	go one.goFuncBatcher()
	go one.goFuncTimer()
	return nil
}
