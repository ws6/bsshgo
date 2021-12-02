package s3io

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

type deadlineConn struct {
	Timeout time.Duration
	net.Conn
}

func (c *deadlineConn) Read(b []byte) (n int, err error) {
	if err = c.Conn.SetDeadline(time.Now().Add(c.Timeout)); err != nil {
		return
	}
	return c.Conn.Read(b)
}

func (c *deadlineConn) Write(b []byte) (n int, err error) {
	if err = c.Conn.SetDeadline(time.Now().Add(c.Timeout)); err != nil {
		return
	}
	return c.Conn.Write(b)
}

const defaultRedirectLimit = 2

// Preserve headers on redirect
// See: https://github.com/golang/go/issues/4800
func redirectPreserveHeaders(req *http.Request, via []*http.Request) error {
	if len(via) == 0 {
		// No redirects
		return nil
	}

	// Only allow GET redirection, prevent possibility of sending data to a malicious server.
	if req.Method != http.MethodGet {
		return fmt.Errorf("method %s redirection not allowed", req.Method)
	}

	if len(via) > defaultRedirectLimit {
		return fmt.Errorf("%d consecutive requests(redirects)", len(via))
	}

	// mutate the subsequent redirect requests with the first Header
	for key, val := range via[0].Header {
		req.Header[key] = val
	}
	return nil
}

// ClientWithTimeout is an http client optimized for high throughput
// to S3, It times out more aggressively than the default
// http client in net/http as well as setting deadlines on the TCP connection.
// keeps original request headers for redirect.
func ClientWithTimeout(timeout time.Duration) *http.Client {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, timeout)
			if err != nil {
				return nil, err
			}
			if tc, ok := c.(*net.TCPConn); ok {
				tc.SetKeepAlive(true)
				tc.SetKeepAlivePeriod(timeout)
			}
			return &deadlineConn{timeout, c}, nil
		},
		ResponseHeaderTimeout: timeout,
		MaxIdleConnsPerHost:   10,
	}
	return &http.Client{Transport: transport, CheckRedirect: redirectPreserveHeaders}
}
