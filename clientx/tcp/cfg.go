package tcp

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/DaiYuANg/arcgo/clientx"
)

// Config configures the TCP client implementation.
type Config struct {
	Network      string
	Address      string
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	KeepAlive    time.Duration
	TLS          clientx.TLSConfig
}

const defaultDialTimeout = 5 * time.Second

// ErrInvalidConfig indicates that the TCP client configuration is invalid.
var ErrInvalidConfig = errors.New("invalid tcp client config")

// NormalizeAndValidate normalizes cfg and validates all supported options.
func (cfg Config) NormalizeAndValidate() (Config, error) {
	out := cfg
	out.Network = strings.TrimSpace(out.Network)
	out.Address = strings.TrimSpace(out.Address)

	if out.Network == "" {
		out.Network = "tcp"
	}
	if out.Address == "" {
		return Config{}, fmt.Errorf("%w: address is required", ErrInvalidConfig)
	}
	if out.DialTimeout == 0 {
		out.DialTimeout = defaultDialTimeout
	}
	if out.DialTimeout < 0 || out.ReadTimeout < 0 || out.WriteTimeout < 0 || out.KeepAlive < 0 {
		return Config{}, fmt.Errorf("%w: timeout values must be >= 0", ErrInvalidConfig)
	}
	if out.TLS.Enabled && !strings.HasPrefix(out.Network, "tcp") {
		return Config{}, fmt.Errorf("%w: tls requires tcp network", ErrInvalidConfig)
	}

	return out, nil
}
