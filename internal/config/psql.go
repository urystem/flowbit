package config

type dbConfig struct {
	host     string
	port     int
	user     string
	password string
	name     string
	sslMode  string
}

func (c *config) initDBConfig() dbConfig {
	dbConf := dbConfig{}
	dbConf.host = mustGetEnvString("DB_HOST")
	dbConf.port = mustGetEnvInt("DB_PORT")
	dbConf.user = mustGetEnvString("DB_USER")
	dbConf.password = mustGetEnvString("DB_PASSWORD")
	dbConf.name = mustGetEnvString("DB_NAME")
	dbConf.sslMode = mustGetEnvString("DB_SSLMODE")
	return dbConf
}

func (dbC *dbConfig) GetHostName() string {
	return dbC.host
}

func (dbC *dbConfig) GetPort() int {
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
