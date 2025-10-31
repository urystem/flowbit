package one

import "context"

type OneMinuteGlobalInter interface {
	Run(ctx context.Context) error
	RedisNotWorking
}

type RedisNotWorking interface {
	IsNotWorking() bool
}
