package config

import (
	"strconv"
	"strings"
	"time"
)

type sources struct {
	hostAdrr map[string]string
	interv   time.Duration
}
type SourcesCfg interface {
	// GetPort(host string) uint16
	// GetCountWorkers() uint8
	GetInterval() time.Duration
	GetAddresses() map[string]string
	// GetCountOfAllWorkers() uint16
}

func (c *config) initSources() sources {
	second := mustGetEnvInt("MARKET_INTERVAL")
	if second < 1 || second > 60 {
		panic("ss")
	}

	addrSlc := mustGetArrayStr("MARKET_ADDRESSES")
	if len(addrSlc) == 0 {
		panic("no addr")
	}
	lenAddr := len(addrSlc)
	checkHost := make(map[string]struct{}, lenAddr)
	checkPort := make(map[uint16]struct{}, lenAddr)
	hostAddr := make(map[string]string)
	for i := range addrSlc {
		address := strings.Split(addrSlc[i], ":")
		if len(address) != 2 {
			panic("")
		}
		hostStr := strings.TrimSpace(address[0])
		portStr := strings.TrimSpace(address[1])
		port, err := strconv.ParseUint(portStr, 10, 16)
		if err != nil {
			panic(err)
		}
		checkHost[hostStr] = struct{}{}
		checkPort[uint16(port)] = struct{}{}
		hostAddr[hostStr] = hostStr + ":" + portStr
	}

	if lenAddr != len(checkHost) {
		panic("duplicated host")
	} else if len(checkPort) < lenAddr {
		panic("duplicated port")
	}
	return sources{
		hostAdrr: hostAddr,
		interv:   time.Duration(second) * time.Second,
	}
}

func (s *sources) GetAddresses() map[string]string { return s.hostAdrr }

func (s *sources) GetInterval() time.Duration { return s.interv }
