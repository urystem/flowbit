package exchange

import (
	"net"
	"strings"
	"sync"

	"marketflow/internal/domain"
	"marketflow/internal/ports/outbound"
)

type stream struct {
	exName string
	con    net.Conn
	using  bool
	mu     sync.Mutex
	get    func() *domain.Exchange
}

func InitStream(addr string, get func() *domain.Exchange) (outbound.StreamAdapterInter, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	// conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", "example.com:1234")
	before, _, _ := strings.Cut(addr, ":")
	return &stream{
		exName: before,
		con:    conn,
		get:    get,
	}, nil
}

func (s *stream) CloseStream() error {
	return s.con.Close()
}
