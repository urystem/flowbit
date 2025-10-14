package config

type redisConfig struct {
	port string
	pass string
}

type RedisConfig interface {
	GetAddr() string
	GetPass() string
}

func (c *config) initRedisConf() redisConfig {
	conf := redisConfig{}
	conf.port = mustGetEnvString("REDIS_PORT")
	conf.pass = mustGetEnvString("REDIS_PASS")
	return conf
}

func (r *redisConfig) GetAddr() string { return r.port }

func (r *redisConfig) GetPass() string { return r.pass }
