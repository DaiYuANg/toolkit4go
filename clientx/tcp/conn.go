package tcp

import (
	"fmt"
	"net"
	"time"

	"github.com/DaiYuANg/arcgo/clientx"
)

type timeoutConn struct {
	net.Conn
	readTimeout  time.Duration
	writeTimeout time.Duration
	addr         string
	hooks        []clientx.Hook
}

func (c *timeoutConn) Read(b []byte) (int, error) {
	return c.runIO("read", b, c.readTimeout, c.SetReadDeadline, c.Conn.Read)
}

func (c *timeoutConn) Write(b []byte) (int, error) {
	return c.runIO("write", b, c.writeTimeout, c.SetWriteDeadline, c.Conn.Write)
}

func (c *timeoutConn) runIO(
	op string,
	data []byte,
	timeout time.Duration,
	setDeadline func(time.Time) error,
	run func([]byte) (int, error),
) (int, error) {
	start := time.Now()
	if err := applyDeadline(setDeadline, timeout, op, c.addr); err != nil {
		emitIO(op, c.addr, 0, time.Since(start), err, c.hooks)
		return 0, err
	}

	n, err := run(data)
	if err != nil {
		wrappedErr := wrapClientError(op, c.addr, err)
		emitIO(op, c.addr, n, time.Since(start), wrappedErr, c.hooks)
		return n, wrappedErr
	}

	emitIO(op, c.addr, n, time.Since(start), nil, c.hooks)
	return n, nil
}

func applyDeadline(setDeadline func(time.Time) error, timeout time.Duration, op, addr string) error {
	if timeout <= 0 {
		return nil
	}
	if err := setDeadline(time.Now().Add(timeout)); err != nil {
		return fmt.Errorf("set tcp %s deadline for %s: %w", op, addr, err)
	}
	return nil
}

func emitIO(op, addr string, bytes int, duration time.Duration, err error, hooks []clientx.Hook) {
	clientx.EmitIO(hooks, clientx.IOEvent{
		Protocol: clientx.ProtocolTCP,
		Op:       op,
		Addr:     addr,
		Bytes:    bytes,
		Duration: duration,
		Err:      err,
	})
}
