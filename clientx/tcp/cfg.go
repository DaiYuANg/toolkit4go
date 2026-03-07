package tcp

import (
	"time"

	"github.com/DaiYuANg/arcgo/clientx"
)

type Config struct {
	Network      string
	Address      string
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	KeepAlive    time.Duration
	TLS          clientx.TLSConfig
}
