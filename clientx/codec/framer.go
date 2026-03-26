package codec

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	// DefaultMaxFrameBytes is the default maximum frame size for LengthPrefixedFramer.
	DefaultMaxFrameBytes uint32 = 4 * 1024 * 1024
)

// Framer reads and writes framed byte payloads.
type Framer interface {
	ReadFrame(r io.Reader) ([]byte, error)
	WriteFrame(w io.Writer, frame []byte) error
}

// LengthPrefixedFramer implements a 4-byte big-endian length-prefixed framing format.
type LengthPrefixedFramer struct {
	MaxFrameBytes uint32
}

// NewLengthPrefixed creates a LengthPrefixedFramer with a validated maximum frame size.
func NewLengthPrefixed(maxFrameBytes uint32) *LengthPrefixedFramer {
	if maxFrameBytes == 0 {
		maxFrameBytes = DefaultMaxFrameBytes
	}
	return &LengthPrefixedFramer{MaxFrameBytes: maxFrameBytes}
}

// ReadFrame reads a single framed payload from r.
func (f *LengthPrefixedFramer) ReadFrame(r io.Reader) ([]byte, error) {
	if r == nil {
		return nil, errors.New("reader is nil")
	}

	var header [4]byte
	if _, err := io.ReadFull(r, header[:]); err != nil {
		return nil, fmt.Errorf("read frame header: %w", err)
	}
	size := binary.BigEndian.Uint32(header[:])
	if size > f.MaxFrameBytes {
		return nil, fmt.Errorf("frame too large: %d > %d", size, f.MaxFrameBytes)
	}
	if size == 0 {
		return []byte{}, nil
	}

	frame := make([]byte, size)
	if _, err := io.ReadFull(r, frame); err != nil {
		return nil, fmt.Errorf("read frame body: %w", err)
	}
	return frame, nil
}

// WriteFrame writes frame to w using the configured length prefix format.
func (f *LengthPrefixedFramer) WriteFrame(w io.Writer, frame []byte) error {
	if w == nil {
		return errors.New("writer is nil")
	}
	frameSize := uint64(len(frame))
	if frameSize > uint64(f.MaxFrameBytes) {
		return fmt.Errorf("frame too large: %d > %d", len(frame), f.MaxFrameBytes)
	}

	var header [4]byte
	//nolint:gosec // frameSize is bounded by MaxFrameBytes, which is a uint32.
	binary.BigEndian.PutUint32(header[:], uint32(frameSize))
	if err := writeFull(w, header[:]); err != nil {
		return err
	}
	return writeFull(w, frame)
}

func writeFull(w io.Writer, data []byte) error {
	for len(data) > 0 {
		n, err := w.Write(data)
		if err != nil {
			return fmt.Errorf("write frame bytes: %w", err)
		}
		data = data[n:]
	}
	return nil
}
