package internal

import (
	"net"
	"time"
)

func isAlive(s *Server) (bool, error) {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", s._url.Host, timeout)
	if err != nil {
		return false, err
	}
	defer conn.Close()
	return true, nil
}
