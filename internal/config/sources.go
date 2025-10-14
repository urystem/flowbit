package config

import (
	"strconv"
	"strings"
	"time"
)

type sources struct {
	address []string
	interv  time.Duration
}
type SourcesCfg interface {
	// GetPort(host string) uint16
	// GetCountWorkers() uint8
	GetInterval() time.Duration
	GetAddresses() []string
	// GetCountOfAllWorkers() uint16
}

func (c *config) initSources() sources {
	second := mustGetEnvInt("MARKET_INTERVAL")
	if second < 1 || second > 60 {
		// slog.Error("ss")
		panic("ss")
	}

	addrSlc := mustGetArrayStr("MARKET_ADDRESSES")
	if len(addrSlc) == 0 {
		panic("no addr")
	}
	lenAddr := len(addrSlc)
	checkHost := make(map[string]struct{}, lenAddr)
	checkPort := make(map[uint16]struct{}, lenAddr)
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
		addrSlc[i] = hostStr + ":" + portStr
	}

	if lenAddr != len(checkHost) {
		panic("duplicated host")
	} else if len(checkPort) < lenAddr {
		panic("duplicated port")
	}
	return sources{
		address: addrSlc,
		interv:  time.Duration(second) * time.Second,
	}
}

func (s *sources) GetAddresses() []string { return s.address }

func (s *sources) GetInterval() time.Duration { return s.interv }
