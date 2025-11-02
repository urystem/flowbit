package one

import "context"

type OneMinuteGlobalInter interface {
	Run(ctx context.Context) error
	RedisNotWorking
	oneInsertBatches
}

type RedisNotWorking interface {
	IsNotWorking() bool
}

type OneForUseCase interface {
	RedisNotWorking
	oneInsertBatches
}

type oneInsertBatches interface {
	PushDone(ctx context.Context)
}
