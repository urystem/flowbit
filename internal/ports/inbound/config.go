package inbound

import "time"

type ServerCfg interface {
	GetPort() uint16
}

type DBConfig interface {
	GetHostName() string
	GetPort() uint16
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
	// GetPort(host string) uint16
	// GetCountWorkers() uint8
	GetInterval() time.Duration
	GetAddresses() []string
	// GetCountOfAllWorkers() uint16
}

type Config interface {
	GetServerCfg() ServerCfg
	GetDBConfig() DBConfig
	GetRedisConfig() RedisConfig
	GetSourcesCfg() SourcesCfg
}
