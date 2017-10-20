package http

import (
	"net"
	"net/http"
	"time"
)

// NewTimeoutClient returns a new http.Client using
// the "timeout" for open connection and read-write operations.
func NewTimeoutClient(timeout time.Duration) *http.Client {
	transport := http.Transport{
		Dial: func(network string, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(network, addr, timeout)
			if err != nil {
				return nil, err
			}
			err = conn.SetDeadline(time.Now().Add(timeout))
			return conn, err
		},
	}

	client := &http.Client{
		Transport: &transport,
	}

	return client
}
