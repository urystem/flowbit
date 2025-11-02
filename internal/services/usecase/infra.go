package usecase

import (
	"marketflow/internal/ports/inbound"
	"marketflow/internal/ports/outbound"
	"marketflow/internal/services/one"
	"marketflow/internal/services/streams"
)

type myUsecase struct {
	strm streams.StreamUsecase  // health
	db   outbound.PgxForUseCase // health
	rdb  outbound.RedisUseCase  // health
	one  one.OneForUseCase
}

func NewUsecase(strm streams.StreamUsecase, db outbound.PgxForUseCase, rdb outbound.RedisUseCase, one one.OneForUseCase) inbound.UsecaseInter {
	return &myUsecase{
		strm: strm,
		db:   db,
		rdb:  rdb,
		one:  one,
	}
}
