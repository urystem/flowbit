package outbound

import "context"

type GeneratorInter interface {
	Check() error
	Start(ctx context.Context)
}
