package config

type sources struct {
	ports        []string
	countWorkers uint8
}

func (c *config) initSources() *sources {
	count := mustGetEnvInt("count")
	if count < 1 || count > 255 {
		// slog.Error("ss")
		panic("ss")
	}
	return &sources{
		ports:        mustGetArrayStr("ss"),
		countWorkers: uint8(count),
	}
}

func (s *sources) GetAddrs() []string {
	return s.ports
}
