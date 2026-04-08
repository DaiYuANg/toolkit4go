// Package randomport provides utilities for finding available random ports.
package randomport

import (
	"context"
	"fmt"
	"net"

	"github.com/DaiYuANg/arcgo/collectionx"
	"github.com/samber/oops"
)

var (
	// usedPorts tracks ports that have been allocated during the current process.
	usedPorts = collectionx.NewConcurrentSet[int]()
)

const maxFindAttempts = 50

// Find returns a random available port that is not currently in use.
// It checks both TCP port availability and tracks previously allocated ports
// to avoid conflicts when multiple servers are started in the same process.
func Find() (int, error) {
	ctx := context.Background()
	var lastErr error

	// Try up to 50 times to find an available port
	for range maxFindAttempts {
		port, err := findAvailablePort(ctx)
		if err != nil {
			lastErr = err
			continue
		}
		if usedPorts.AddIfAbsent(port) {
			return port, nil
		}
	}

	if lastErr != nil {
		return 0, oops.In("pkg/randomport").
			With("op", "find_port", "attempts", maxFindAttempts).
			Wrapf(lastErr, "find available port")
	}

	return 0, oops.In("pkg/randomport").
		With("op", "find_port", "attempts", maxFindAttempts).
		New("failed to find available port")
}

// findAvailablePort finds a single available port by listening on port 0.
func findAvailablePort(ctx context.Context) (port int, err error) {
	var listenConfig net.ListenConfig
	listener, err := listenConfig.Listen(ctx, "tcp", "127.0.0.1:0")
	if err != nil {
		return 0, oops.In("pkg/randomport").
			With("op", "listen_port", "network", "tcp", "addr", "127.0.0.1:0").
			Wrapf(err, "listen for random port")
	}
	defer func() {
		closeErr := listener.Close()
		if err == nil && closeErr != nil {
			err = oops.In("pkg/randomport").
				With("op", "close_listener", "network", "tcp", "addr", listener.Addr().String()).
				Wrapf(closeErr, "close random port listener")
		}
	}()

	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		return 0, oops.In("pkg/randomport").
			With("op", "resolve_port", "addr_type", fmt.Sprintf("%T", listener.Addr())).
			Errorf("unexpected listener address type %T", listener.Addr())
	}

	return addr.Port, nil
}

// Release releases a port back to the available pool.
// This is primarily useful for testing scenarios.
func Release(port int) {
	usedPorts.Remove(port)
}

// MustFind returns a random available port or panics if none can be found.
func MustFind() int {
	port, err := Find()
	if err != nil {
		panic(oops.In("pkg/randomport").
			With("op", "must_find_port").
			Wrapf(err, "find available port"))
	}

	return port
}
