package one

import "context"

type OneMinuteGlobalInter interface {
	Run(ctx context.Context) error
	OneMinuteStatus
}

type OneMinuteStatus interface {
	IsNotWorking() bool
}
