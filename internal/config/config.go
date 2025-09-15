package config

import "marketflow/internal/ports/inbound"

type config struct {
	server serverCfg

	db dbConfig

	redis redisConfig
	src   sources
	wr    workerCfg
}

func Load() inbound.Config {
	conf := &config{}
	// conf.server = conf.initServerCfg()
	// conf.db = conf.initDBConfig()
	conf.redis = conf.initRedisConf()
	conf.src = conf.initSources()
	conf.wr = conf.initWorkerPoolCfg()

	return conf
}

func (conf *config) GetServerCfg() inbound.ServerCfg { return &conf.server }

func (conf *config) GetDBConfig() inbound.DBConfig { return &conf.db }

func (conf *config) GetRedisConfig() inbound.RedisConfig { return &conf.redis }

func (conf *config) GetSourcesCfg() inbound.SourcesCfg { return &conf.src }

func (conf *config) GetWorkerCfg() inbound.WorkerCfg { return &conf.wr }
