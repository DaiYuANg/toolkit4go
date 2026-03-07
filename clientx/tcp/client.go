package tcp

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
)

type Client struct {
	cfg Config
}

func New(cfg Config, opts ...Option) *Client {
	c := &Client{cfg: cfg}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) Dial(ctx context.Context) (net.Conn, error) {
	network := c.cfg.Network
	if network == "" {
		network = "tcp"
	}

	dialer := &net.Dialer{
		Timeout:   c.cfg.DialTimeout,
		KeepAlive: c.cfg.KeepAlive,
	}

	if c.cfg.TLS.Enabled {
		conn, err := tls.DialWithDialer(dialer, network, c.cfg.Address, &tls.Config{
			InsecureSkipVerify: c.cfg.TLS.InsecureSkipVerify,
			ServerName:         c.cfg.TLS.ServerName,
		})
		if err != nil {
			return nil, fmt.Errorf("dial tls tcp failed: %w", err)
		}
		defer func(conn *tls.Conn) {
			_ = conn.Close()
		}(conn)
		return wrapConn(conn, c.cfg), nil
	}

	conn, err := dialer.DialContext(ctx, network, c.cfg.Address)
	if err != nil {
		return nil, fmt.Errorf("dial tcp failed: %w", err)
	}
	return wrapConn(conn, c.cfg), nil
}

func wrapConn(conn net.Conn, cfg Config) net.Conn {
	return &timeoutConn{
		Conn:         conn,
		readTimeout:  cfg.ReadTimeout,
		writeTimeout: cfg.WriteTimeout,
	}
}
