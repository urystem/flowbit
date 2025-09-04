package config

import (
	"strconv"
	"strings"
)

type sources struct {
	address      []string
	countWorkers uint8
}

func (c *config) initSources() sources {
	count := mustGetEnvInt("MARKET_DEFAULT_WORKERS")
	if count < 1 || count > 255 {
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
		addrSlc[i] = hostStr + portStr
	}

	if lenAddr != len(checkHost) {
		panic("duplicated host")
	} else if len(checkPort) < lenAddr {
		panic("duplicated port")
	}
	return sources{
		address:      addrSlc,
		countWorkers: uint8(count),
	}
}

func (s *sources) GetHosts() []string { return s.address }

func (s *sources) GetCountWorkers() uint8 { return s.countWorkers }
