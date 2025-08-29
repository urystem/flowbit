package source

import "sync"

type Market struct {
	// test  bool
	addrs        []string
	countWorkers uint8
	wg           sync.WaitGroup
}

func InitSource()
