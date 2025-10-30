package config

type config struct {
	server serverCfg
	db     dbConfig
	redis  redisConfig
	src    sources
	wr     workerCfg
}

type ConfigInter interface {
	GetServerCfg() ServerCfg
	GetDBConfig() DBConfig
	GetRedisConfig() RedisConfig
	GetSourcesCfg() SourcesCfg
	GetWorkerCfg() WorkerCfg
}

func Load() ConfigInter {
	conf := &config{}
	conf.server = conf.initServerCfg()
	conf.db = conf.initDBConfig()
	conf.redis = conf.initRedisConf()
	conf.src = conf.initSources()
	conf.wr = conf.initWorkerPoolCfg()

	return conf
}

func (conf *config) GetServerCfg() ServerCfg { return &conf.server }

func (conf *config) GetDBConfig() DBConfig { return &conf.db }

func (conf *config) GetRedisConfig() RedisConfig { return &conf.redis }

func (conf *config) GetSourcesCfg() SourcesCfg { return &conf.src }

func (conf *config) GetWorkerCfg() WorkerCfg { return &conf.wr }
