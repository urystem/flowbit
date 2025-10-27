package config

import (
	"os"
	"strings"
	"time"
)

type sources struct {
	testAddr string
	hostAdrr map[string]string
	interv   time.Duration
}
type SourcesCfg interface {
	GetTestAddr() string
	GetInterval() time.Duration
	GetAddresses() map[string]string
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

	checkAddr := make(map[string]struct{})
	testAddr := mustGetEnvString("MARKET_TEST_ADDRESS")
	checkAddr[testAddr] = struct{}{}

	exchanges := make(map[string]string)
	for _, v := range addrSlc {
		checkAddr[v] = struct{}{}
		address := strings.Split(v, ":")
		if len(address) != 2 {
			os.Exit(1)
		}
		exchanges[address[0]] = v
	}
	if len(checkAddr)-1 != len(addrSlc) {
		os.Exit(1)
	}
	return sources{
		testAddr: testAddr,
		hostAdrr: exchanges,
		interv:   time.Duration(second) * time.Second,
	}
}

func (s *sources) GetAddresses() map[string]string { return s.hostAdrr }

func (s *sources) GetInterval() time.Duration { return s.interv }

func (s *sources) GetTestAddr() string { return s.testAddr }
