package http

import (
	"time"

	"github.com/DaiYuANg/arcgo/clientx"
	"github.com/DaiYuANg/arcgo/collectionx"
)

type Config struct {
	BaseURL   string
	Timeout   time.Duration
	Headers   collectionx.Map[string, string]
	UserAgent string
	Retry     clientx.RetryConfig
	TLS       clientx.TLSConfig
}
