package websocket

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/lxzan/gws"
)

type gwsConn struct {
	socket *gws.Conn
	opts   Options
	recv   chan Message
	errCh  chan error
	once   sync.Once
}

type eventBridge struct {
	conn *gwsConn
}

// HandlerFunc adapts a WebSocket handler to net/http.
func HandlerFunc(handler Handler, options ...Option) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := Upgrade(w, r, handler, options...); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}

// Upgrade upgrades the request to a WebSocket connection and runs handler.
func Upgrade(w http.ResponseWriter, r *http.Request, handler Handler, options ...Option) error {
	if handler == nil {
		return fmt.Errorf("%w: nil handler", ErrUpgradeFailed)
	}
	cfg := applyOptions(options)
	bridgeConn := &gwsConn{
		opts:  cfg,
		recv:  make(chan Message, 32),
		errCh: make(chan error, 1),
	}
	upgrader := gws.NewUpgrader(&eventBridge{conn: bridgeConn}, &gws.ServerOption{
		HandshakeTimeout:   cfg.HandshakeTimeout,
		ReadMaxPayloadSize: cfg.MaxMessageSize,
		PermessageDeflate:  gws.PermessageDeflate{Enabled: cfg.EnableCompression},
		Authorize: func(req *http.Request, _ gws.SessionStorage) bool {
			if cfg.CheckOrigin == nil {
				return true
			}
			return cfg.CheckOrigin(req)
		},
	})
	socket, err := upgrader.Upgrade(w, r)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrUpgradeFailed, err)
	}
	bridgeConn.socket = socket
	go socket.ReadLoop()

	ctx := r.Context()
	if runErr := handler(ctx, bridgeConn); runErr != nil {
		closeErr := bridgeConn.Close(1011, []byte(runErr.Error()))
		if closeErr != nil {
			return errors.Join(runErr, fmt.Errorf("httpx/websocket: close connection after handler error: %w", closeErr))
		}
		return runErr
	}
	return nil
}

// Read reads the next WebSocket message or returns when ctx is done.
func (c *gwsConn) Read(ctx Context) (Message, error) {
	if ctx == nil {
		return Message{}, errors.New("httpx/websocket: nil read context")
	}
	select {
	case msg, ok := <-c.recv:
		if !ok {
			return Message{}, ErrClosed
		}
		return msg, nil
	case err := <-c.errCh:
		if err == nil {
			return Message{}, ErrClosed
		}
		return Message{}, err
	case <-ctx.Done():
		return Message{}, fmt.Errorf("httpx/websocket: read canceled: %w", ctx.Err())
	}
}

// Write writes a typed WebSocket message to the connection.
func (c *gwsConn) Write(msg Message) (retErr error) {
	if c.socket == nil {
		return ErrClosed
	}
	cleanup, err := c.startWrite()
	if err != nil {
		return err
	}
	defer func() {
		if err := cleanup(); err != nil {
			if retErr == nil {
				retErr = err
				return
			}
			retErr = errors.Join(retErr, err)
		}
	}()

	return c.writeFrame(msg)
}

func (c *gwsConn) startWrite() (func() error, error) {
	if c.opts.WriteTimeout <= 0 {
		return noopWriteCleanup, nil
	}

	if err := c.socket.SetWriteDeadline(time.Now().Add(c.opts.WriteTimeout)); err != nil {
		return nil, fmt.Errorf("httpx/websocket: set write deadline: %w", err)
	}

	return func() error {
		if err := c.socket.SetWriteDeadline(time.Time{}); err != nil {
			return fmt.Errorf("httpx/websocket: reset write deadline: %w", err)
		}
		return nil
	}, nil
}

func noopWriteCleanup() error {
	return nil
}

func (c *gwsConn) writeFrame(msg Message) error {
	opcode, ok := toGWSOpcode(msg.Type)
	if !ok {
		return fmt.Errorf("httpx/websocket: unsupported message type: %d", msg.Type)
	}
	if msg.Type == MessageClose {
		if err := c.socket.WriteClose(1000, msg.Data); err != nil {
			return fmt.Errorf("httpx/websocket: write close frame: %w", err)
		}
		return nil
	}
	if err := c.socket.WriteMessage(opcode, msg.Data); err != nil {
		return fmt.Errorf("httpx/websocket: write message: %w", err)
	}
	return nil
}

// Close sends a WebSocket close frame.
func (c *gwsConn) Close(code uint16, reason []byte) error {
	if c.socket == nil {
		return nil
	}
	err := c.socket.WriteClose(code, reason)
	if err != nil && !errors.Is(err, gws.ErrConnClosed) {
		return fmt.Errorf("httpx/websocket: write close frame: %w", err)
	}
	return nil
}

func (b *eventBridge) OnOpen(socket *gws.Conn) {
	b.refreshReadDeadlines(socket)
}

func (b *eventBridge) OnClose(_ *gws.Conn, err error) {
	b.conn.once.Do(func() {
		if err != nil {
			b.conn.errCh <- err
		}
		close(b.conn.recv)
		close(b.conn.errCh)
	})
}

func (b *eventBridge) OnPing(socket *gws.Conn, payload []byte) {
	b.refreshReadDeadlines(socket)
	if err := socket.WritePong(payload); err != nil {
		b.conn.reportError(fmt.Errorf("httpx/websocket: write pong frame: %w", err))
	}
}

func (b *eventBridge) OnPong(socket *gws.Conn, _ []byte) {
	b.refreshIdleDeadline(socket)
}

func (b *eventBridge) OnMessage(socket *gws.Conn, message *gws.Message) {
	defer func() {
		if err := message.Close(); err != nil {
			b.conn.reportError(fmt.Errorf("httpx/websocket: close message: %w", err))
		}
	}()
	msgType, ok := fromGWSOpcode(message.Opcode)
	if !ok {
		return
	}
	b.refreshReadDeadlines(socket)
	payload := append([]byte(nil), message.Bytes()...)
	b.conn.recv <- Message{Type: msgType, Data: payload}
}

func toGWSOpcode(t MessageType) (gws.Opcode, bool) {
	switch t {
	case MessageText:
		return gws.OpcodeText, true
	case MessageBinary:
		return gws.OpcodeBinary, true
	case MessagePing:
		return gws.OpcodePing, true
	case MessagePong:
		return gws.OpcodePong, true
	case MessageClose:
		return gws.OpcodeCloseConnection, true
	default:
		return 0, false
	}
}

func fromGWSOpcode(op gws.Opcode) (MessageType, bool) {
	switch op {
	case gws.OpcodeText:
		return MessageText, true
	case gws.OpcodeBinary:
		return MessageBinary, true
	case gws.OpcodePing:
		return MessagePing, true
	case gws.OpcodePong:
		return MessagePong, true
	case gws.OpcodeCloseConnection:
		return MessageClose, true
	case gws.OpcodeContinuation:
		return 0, false
	default:
		return 0, false
	}
}

func (c *gwsConn) reportError(err error) {
	if err == nil {
		return
	}

	defer func() {
		if recovered := recover(); recovered != nil {
			return
		}
	}()

	select {
	case c.errCh <- err:
	default:
	}
}

func (b *eventBridge) refreshReadDeadlines(socket *gws.Conn) {
	b.refreshIdleDeadline(socket)
	if b.conn.opts.ReadTimeout > 0 {
		if err := socket.SetReadDeadline(time.Now().Add(b.conn.opts.ReadTimeout)); err != nil {
			b.conn.reportError(fmt.Errorf("httpx/websocket: set read deadline: %w", err))
		}
	}
}

func (b *eventBridge) refreshIdleDeadline(socket *gws.Conn) {
	if b.conn.opts.IdleTimeout > 0 {
		if err := socket.SetDeadline(time.Now().Add(b.conn.opts.IdleTimeout)); err != nil {
			b.conn.reportError(fmt.Errorf("httpx/websocket: set idle deadline: %w", err))
		}
	}
}
