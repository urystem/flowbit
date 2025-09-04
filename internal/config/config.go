package config

import "marketflow/internal/ports/inbound"

type config struct {
	server serverCfg

	db dbConfig

	redis redisConfig
	src   sources
}

func Load() inbound.Config {
	conf := &config{}
	// conf.server = conf.initServerCfg()
	// conf.db = conf.initDBConfig()
	// conf.redis = conf.initRedisConf()
	conf.src = conf.initSources()

	return conf
}

func (conf *config) GetServerCfg() inbound.ServerCfg {
	return &conf.server
}

func (conf *config) GetDBConfig() inbound.DBConfig {
	return &conf.db
}

func (conf *config) GetRedisConfig() inbound.RedisConfig {
	return &conf.redis
}

func (conf *config) GetSourcesCfg() inbound.SourcesCfg {
	return &conf.src
}
