package usecase

func (u *myUsecase) SwitchToTest() error {
	u.strm.StopJustStreams()
	return u.strm.StartTestStream()
}

func (u *myUsecase) SwitchToLive() {
	u.strm.StopTestStream()
	u.strm.StartJustStreams()
}
