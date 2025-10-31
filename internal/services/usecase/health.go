package usecase

import "context"

type health struct {
	Streams map[string]string `json:"streams"`
	DB      string            `json:"pgx"`
	RDB     string            `json:"redis"`
}

func (u *myUsecase) CheckHealth(ctx context.Context) any {
	h := &health{
		Streams: map[string]string{},
	}

	for ex, err := range u.strm.CheckHealth() {
		if err != nil {
			h.Streams[ex] = err.Error()
		} else {
			h.Streams[ex] = "ok"
		}
	}
	err := u.db.CheckHealth(ctx)
	if err != nil {
		h.DB = err.Error()
	} else {
		h.DB = "ok"
	}

	err = u.rdb.CheckHealth(ctx)
	if err != nil {
		h.RDB = err.Error()
	} else {
		h.RDB = "ok"
	}
	return h
}
