package config

import "time"

type workerCfg struct {
	maxOrDefWorker int
	elastic        bool
	interv         time.Duration
}

type WorkerCfg interface {
	GetBoolElasticWorker() bool
	GetCountOfMaxOrDefWorker() int
	GetElasticInterval() time.Duration
}

func (c *config) initWorkerPoolCfg() workerCfg {
	count := mustGetEnvInt("MARKET_DEFAULT_WORKERS")
	if count < 0 || count > 65000 {
		panic("worker count must beeen between 0 and uint16")
	}

	intv := mustGetEnvInt("MARKET_ELASTIC_INTERVAL")
	if intv < 0 || intv > 60 {
		panic("ss")
	}
	return workerCfg{
		maxOrDefWorker: count,
		elastic:        mustGetBoolean("MARKET_ELASTIC_WORKER_POOL"),
		interv:         time.Duration(intv * int(time.Second)),
	}
}

func (w *workerCfg) GetBoolElasticWorker() bool { return w.elastic }

func (w *workerCfg) GetCountOfMaxOrDefWorker() int { return w.maxOrDefWorker }

func (w *workerCfg) GetElasticInterval() time.Duration { return w.interv }
