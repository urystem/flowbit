package config

import (
	"strconv"
	"strings"
)

type hosts struct {
	ex         map[string]market
	allWorkers uint16
}

type market struct {
	port         uint16
	CountWorkers uint8
}

func (c *config) initSources() hosts {
	count := mustGetEnvInt("MARKET_DEFAULT_WORKERS")
	if count < 1 || count > 255 {
		// slog.Error("ss")
		panic("ss")
	}
	defCount := uint8(count)
	addrSlc := mustGetArrayStr("MARKET_ADDRESSES")
	if len(addrSlc) == 0 {
		panic("no addr")
	}

	myHost := make(map[string]market)
	checkPort := make(map[uint16]struct{})
	var workerCount uint16
	for _, addr := range addrSlc {
		host0 := strings.Split(addr, ":")
		if len(host0) != 2 {
			panic("")
		}

		var m market
		portAndWorker := strings.Split(host0[1], "*")
		port, err := strconv.ParseUint(strings.TrimSpace(portAndWorker[0]), 10, 16)
		if err != nil {
			panic(err)
		}

		if len(portAndWorker) > 2 {
			panic("ddd")
		}
		m.port = uint16(port)
		checkPort[m.port] = struct{}{}
		if len(portAndWorker) == 2 {
			workerCount, err := strconv.ParseUint(strings.TrimSpace(portAndWorker[1]), 10, 8)
			if err != nil {
				panic(err)
			}
			m.CountWorkers = uint8(workerCount)
		} else {
			m.CountWorkers = defCount
		}
		if workerCount+uint16(m.CountWorkers) > workerCount {
			workerCount += uint16(m.CountWorkers)
		} else {
			panic("too many workers")
		}
		myHost[strings.TrimSpace(host0[0])] = m
	}

	if len(addrSlc) != len(myHost) {
		panic("duplicated address")
	} else if len(checkPort) < len(myHost) {
		panic("duplicated port")
	} else if len(checkPort) > len(myHost) {
		panic("duplicated host")
	}
	return hosts{
		ex:         myHost,
		allWorkers: workerCount,
	}
}

func (s *hosts) GetHosts() []string {
	hosts := make([]string, 0, len(s.ex))
	for k := range s.ex {
		hosts = append(hosts, k)
	}
	return hosts
}

func (h *hosts) GetCountWorkers(hostKey string) uint8 {
	return h.ex[hostKey].CountWorkers
}

func (h *hosts) GetPort(host string) uint16 {
	return h.ex[host].port
}

func (h *hosts) GetCountOfAllWorkers() uint16 {
	return h.allWorkers
}
