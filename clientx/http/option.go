package http

import "resty.dev/v3"

type Option func(*Client)

func WithRequestMiddleware(fn func(*resty.Client, *resty.Request) error) Option {
	return func(c *Client) {
		c.raw.AddRequestMiddleware(fn)
	}
}

func WithResponseMiddleware(fn func(*resty.Client, *resty.Response) error) Option {
	return func(c *Client) {
		c.raw.AddResponseMiddleware(fn)
	}
}

func WithHeader(key, value string) Option {
	return func(c *Client) {
		c.raw.SetHeader(key, value)
	}
}
