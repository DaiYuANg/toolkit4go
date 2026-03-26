package udp

import (
	"errors"
	"fmt"
	"net"

	clientcodec "github.com/DaiYuANg/arcgo/clientx/codec"
)

const maxUDPPacketSize = 64 * 1024

// CodecConn wraps a connected UDP socket with codec helpers.
type CodecConn struct {
	conn  net.Conn
	codec clientcodec.Codec
	addr  string
}

// NewCodecConn wraps conn with codec helpers.
func NewCodecConn(conn net.Conn, codec clientcodec.Codec, addr string) *CodecConn {
	return &CodecConn{
		conn:  conn,
		codec: codec,
		addr:  addr,
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
		return fmt.Errorf("close udp codec conn: %w", err)
	}
	return nil
}

// WriteValue encodes v and writes one UDP payload.
func (c *CodecConn) WriteValue(v any) error {
	if c.codec == nil {
		return wrapCodecError("encode", c.addr, errors.New("codec is nil"))
	}
	payload, err := c.codec.Marshal(v)
	if err != nil {
		return wrapCodecError("encode", c.addr, err)
	}
	if _, err := c.conn.Write(payload); err != nil {
		return wrapClientError("write", c.addr, err)
	}
	return nil
}

// ReadValue reads one UDP payload and decodes it into v.
func (c *CodecConn) ReadValue(v any) error {
	if c.codec == nil {
		return wrapCodecError("decode", c.addr, errors.New("codec is nil"))
	}

	buf := make([]byte, maxUDPPacketSize)
	n, err := c.conn.Read(buf)
	if err != nil {
		return wrapClientError("read", c.addr, err)
	}
	if err := c.codec.Unmarshal(buf[:n], v); err != nil {
		return wrapCodecError("decode", c.addr, err)
	}
	return nil
}

// CodecPacketConn wraps a packet listener with codec helpers.
type CodecPacketConn struct {
	conn  net.PacketConn
	codec clientcodec.Codec
	addr  string
}

// NewCodecPacketConn wraps conn with packet codec helpers.
func NewCodecPacketConn(conn net.PacketConn, codec clientcodec.Codec, addr string) *CodecPacketConn {
	return &CodecPacketConn{
		conn:  conn,
		codec: codec,
		addr:  addr,
	}
}

// Raw returns the underlying net.PacketConn.
func (c *CodecPacketConn) Raw() net.PacketConn {
	return c.conn
}

// Close closes the underlying packet connection.
func (c *CodecPacketConn) Close() error {
	if c.conn == nil {
		return nil
	}
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("close udp packet codec conn: %w", err)
	}
	return nil
}

// ReadValueFrom reads one packet and decodes it into v.
func (c *CodecPacketConn) ReadValueFrom(v any) (net.Addr, error) {
	if c.codec == nil {
		return nil, wrapCodecError("decode", c.addr, errors.New("codec is nil"))
	}

	buf := make([]byte, maxUDPPacketSize)
	n, addr, err := c.conn.ReadFrom(buf)
	if err != nil {
		return nil, wrapClientError("read_from", c.addr, err)
	}
	if err := c.codec.Unmarshal(buf[:n], v); err != nil {
		return nil, wrapCodecError("decode", c.addr, err)
	}
	return addr, nil
}

// WriteValueTo encodes v and writes it to addr.
func (c *CodecPacketConn) WriteValueTo(v any, addr net.Addr) error {
	if c.codec == nil {
		return wrapCodecError("encode", c.addr, errors.New("codec is nil"))
	}
	payload, err := c.codec.Marshal(v)
	if err != nil {
		return wrapCodecError("encode", c.addr, err)
	}

	if _, err := c.conn.WriteTo(payload, addr); err != nil {
		target := c.addr
		if addr != nil {
			target = addr.String()
		}
		return wrapClientError("write_to", target, err)
	}
	return nil
}
