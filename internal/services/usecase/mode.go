package usecase

func (u *myUsecase) SwitchToTest()  {
	u.strm.StopJustStreams()
	u.strm.StartTestStream()
}

func (u *myUsecase) SwitchToLive() {
	u.strm.StopTestStream()
	u.strm.StartJustStreams()
}
