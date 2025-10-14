package config

import (
	"math"
)

type serverCfg struct {
	port uint16
}

type ServerCfg interface {
	GetPort() uint16
}

func (c *config) initServerCfg() serverCfg {
	portTemp := mustGetEnvInt("SERVER_PORT")
	if portTemp < 0 || math.MaxUint16 < portTemp {
		panic("")
	}
	return serverCfg{port: uint16(portTemp)}
}

func (scfg *serverCfg) GetPort() uint16 { return scfg.port }
