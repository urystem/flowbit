package usecase

import "marketflow/internal/services/streams"

type myUsecase struct {
	rdb  any                   // health
	db   any                   // health
	strm streams.StreamUsecase // health
}

func NewUsecase(strm streams.StreamUsecase) any {
	return nil
}
