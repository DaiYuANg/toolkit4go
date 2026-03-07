package http

import (
	"crypto/tls"
	"net/http"

	"github.com/samber/lo"
	"resty.dev/v3"
)

type Client struct {
	raw *resty.Client
}

func New(cfg Config, opts ...Option) *Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: cfg.TLS.InsecureSkipVerify,
			ServerName:         cfg.TLS.ServerName,
		},
	}

	c := resty.New().
		SetBaseURL(cfg.BaseURL).
		SetTimeout(cfg.Timeout).
		SetTransport(transport)

	if cfg.UserAgent != "" {
		c.SetHeader("User-Agent", cfg.UserAgent)
	}
	for k, v := range cfg.Headers.All() {
		c.SetHeader(k, v)
	}

	if cfg.Retry.Enabled {
		c.SetRetryCount(cfg.Retry.MaxRetries)
		c.SetRetryWaitTime(cfg.Retry.WaitMin)
		c.SetRetryMaxWaitTime(cfg.Retry.WaitMax)
	}

	client := &Client{raw: c}
	lo.ForEach(opts, func(opt Option, index int) {
		opt(client)
	})
	return client
}

func (c *Client) Raw() *resty.Client {
	return c.raw
}

func (c *Client) R() *resty.Request {
	return c.raw.R()
}
