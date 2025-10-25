package exchange

import (
	"net"
	"time"
)

func (s *stream) PingStream() (string, error) {
	// conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", "example.com:1234")
	conn, err := net.DialTimeout("tcp", s.addr, 3*time.Second)
	if err != nil {
		return s.exName, err
	}
	return s.exName, conn.Close()
}
