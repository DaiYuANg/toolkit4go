package tcp

import (
	"errors"
	"fmt"
	"net"

	clientcodec "github.com/DaiYuANg/arcgo/clientx/codec"
)

// CodecConn wraps a TCP connection with codec and framer helpers.
type CodecConn struct {
	conn   net.Conn
	codec  clientcodec.Codec
	framer clientcodec.Framer
	addr   string
}

// NewCodecConn wraps conn with codec/framer helpers.
func NewCodecConn(conn net.Conn, codec clientcodec.Codec, framer clientcodec.Framer, addr string) *CodecConn {
	return &CodecConn{
		conn:   conn,
		codec:  codec,
		framer: framer,
		addr:   addr,
	}
}

// Raw returns the underlying net.Conn.
func (c *CodecConn) Raw() net.Conn {
	return c.conn
}

// Close closes the underlying connection.
func (c *CodecConn) Close() error {
	if c.conn == nil {
		return nil
	}
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("close tcp codec conn: %w", err)
	}
	return nil
}

// WriteValue encodes v and writes it as one framed payload.
func (c *CodecConn) WriteValue(v any) error {
	if c.codec == nil || c.framer == nil {
		return wrapCodecError("encode", c.addr, errors.New("codec/framer is nil"))
	}

	payload, err := c.codec.Marshal(v)
	if err != nil {
		return wrapCodecError("encode", c.addr, err)
	}
	if err := c.framer.WriteFrame(c.conn, payload); err != nil {
		return wrapClientError("write_frame", c.addr, err)
	}
	return nil
}

// ReadValue reads one framed payload and decodes it into v.
func (c *CodecConn) ReadValue(v any) error {
	if c.codec == nil || c.framer == nil {
		return wrapCodecError("decode", c.addr, errors.New("codec/framer is nil"))
	}

	frame, err := c.framer.ReadFrame(c.conn)
	if err != nil {
		return wrapClientError("read_frame", c.addr, err)
	}
	if err := c.codec.Unmarshal(frame, v); err != nil {
		return wrapCodecError("decode", c.addr, err)
	}
	return nil
}
