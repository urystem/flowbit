package config

type serverCfg struct {
	port int
}

func (c *config) initServerCfg() serverCfg {
	cfg := serverCfg{}
	cfg.port = mustGetEnvInt("SERVER_PORT")
	return cfg
}

func (scfg *serverCfg) GetPort() int {
	return scfg.port
}
