package app

import (
	"marketflow/internal/adapters/driven/exchange"
	"marketflow/internal/config"
	"marketflow/internal/ports/outbound"
	"marketflow/internal/services/streams"
	syncpool "marketflow/internal/services/syncPool"
)

func (app *myApp) initStreamsService(cfg config.SourcesCfg, getter syncpool.Getter) (streams.StreamsInter, error) {
	inter := cfg.GetInterval()
	addrMap := cfg.GetAddresses()
	strms := make([]outbound.StreamAdapterInter, 0, len(addrMap))
	for name, addr := range addrMap {
		strm := exchange.InitStream(name, addr, inter, getter.GetNewExchange)
		err := strm.PingStream()
		if err != nil {
			return nil, err
		}
		strms = append(strms, strm)
	}
	return streams.InitStreams(strms, getter), nil
}
