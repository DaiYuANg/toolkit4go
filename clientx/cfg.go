package clientx

import "time"

type RetryConfig struct {
	Enabled    bool
	MaxRetries int
	WaitMin    time.Duration
	WaitMax    time.Duration
}

type TLSConfig struct {
	Enabled            bool
	InsecureSkipVerify bool
	ServerName         string
}
