package config

import "math"

type dbConfig struct {
	host     string
	port     uint16
	user     string
	password string
	name     string
	sslMode  string
}

func (c *config) initDBConfig() dbConfig {
	dbConf := dbConfig{}
	dbConf.host = mustGetEnvString("DB_HOST")
	temPort := mustGetEnvInt("DB_PORT")
	if temPort < 0 || math.MaxUint16 < temPort {
		panic("")
	}
	dbConf.user = mustGetEnvString("DB_USER")
	dbConf.password = mustGetEnvString("DB_PASSWORD")
	dbConf.name = mustGetEnvString("DB_NAME")
	dbConf.sslMode = mustGetEnvString("DB_SSLMODE")
	return dbConf
}

func (dbC *dbConfig) GetHostName() string {
	return dbC.host
}

func (dbC *dbConfig) GetPort() uint16 {
	return dbC.port
}

func (dbC *dbConfig) GetUser() string {
	return dbC.user
}

func (dbC *dbConfig) GetPassword() string {
	return dbC.password
}

func (dbC *dbConfig) GetDBName() string {
	return dbC.name
}

func (dbC *dbConfig) GetSSLMode() string {
	return dbC.sslMode
}
