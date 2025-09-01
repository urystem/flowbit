package inbound

type ServerCfg interface {
	GetPort() int
}

type DBConfig interface {
	GetHostName() string
	GetPort() int
	GetUser() string
	GetPassword() string
	GetDBName() string
	GetSSLMode() string
}

type RedisConfig interface {
	GetAddr() string
	GetPass() string
}

type SourcesCfg interface {
	GetPort(host string) uint16
	GetCountWorkers(hostKey string) uint8
	GetHosts() []string
	GetCountOfAllWorkers() uint16
}

type Config interface {
	GetServerCfg() ServerCfg
	GetDBConfig() DBConfig
	GetRedisConfig() RedisConfig
	GetSourcesCfg() SourcesCfg
}
