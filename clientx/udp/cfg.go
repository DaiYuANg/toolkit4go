package udp

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Config configures the UDP client implementation.
type Config struct {
	Network      string
	Address      string
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

const defaultDialTimeout = 5 * time.Second

// ErrInvalidConfig indicates that the UDP client configuration is invalid.
var ErrInvalidConfig = errors.New("invalid udp client config")

// NormalizeAndValidate normalizes cfg and validates all supported options.
func (cfg Config) NormalizeAndValidate() (Config, error) {
	out := cfg
	out.Network = strings.TrimSpace(out.Network)
	out.Address = strings.TrimSpace(out.Address)

	if out.Network == "" {
		out.Network = "udp"
	}
	if out.Address == "" {
		return Config{}, fmt.Errorf("%w: address is required", ErrInvalidConfig)
	}
	if out.DialTimeout == 0 {
		out.DialTimeout = defaultDialTimeout
	}
	if out.DialTimeout < 0 || out.ReadTimeout < 0 || out.WriteTimeout < 0 {
		return Config{}, fmt.Errorf("%w: timeout values must be >= 0", ErrInvalidConfig)
	}
	if !strings.HasPrefix(out.Network, "udp") {
		return Config{}, fmt.Errorf("%w: network must be udp/udp4/udp6", ErrInvalidConfig)
	}

	return out, nil
}
