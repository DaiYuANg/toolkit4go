package clientx

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"syscall"

	"github.com/samber/lo"
)

// Protocol identifies the transport used by a client operation.
type Protocol string

const (
	// ProtocolUnknown indicates that the transport is not known.
	ProtocolUnknown Protocol = "unknown"
	// ProtocolHTTP identifies HTTP client operations.
	ProtocolHTTP Protocol = "http"
	// ProtocolTCP identifies TCP client operations.
	ProtocolTCP Protocol = "tcp"
	// ProtocolUDP identifies UDP client operations.
	ProtocolUDP Protocol = "udp"
)

// ErrorKind classifies client errors into portable categories.
type ErrorKind string

const (
	// ErrorKindUnknown indicates that the error could not be classified.
	ErrorKindUnknown ErrorKind = "unknown"
	// ErrorKindCanceled indicates that the operation was canceled.
	ErrorKindCanceled ErrorKind = "canceled"
	// ErrorKindTimeout indicates that the operation exceeded a deadline.
	ErrorKindTimeout ErrorKind = "timeout"
	// ErrorKindTemporary indicates that the error may be transient.
	ErrorKindTemporary ErrorKind = "temporary"
	// ErrorKindConnRefused indicates that the remote endpoint refused the connection.
	ErrorKindConnRefused ErrorKind = "conn_refused"
	// ErrorKindDNS indicates that DNS resolution failed.
	ErrorKindDNS ErrorKind = "dns"
	// ErrorKindTLS indicates that TLS negotiation or certificate validation failed.
	ErrorKindTLS ErrorKind = "tls"
	// ErrorKindClosed indicates that the connection or resource is already closed.
	ErrorKindClosed ErrorKind = "closed"
	// ErrorKindNetwork indicates a generic network transport failure.
	ErrorKindNetwork ErrorKind = "network"
	// ErrorKindCodec indicates encoding or decoding failures.
	ErrorKindCodec ErrorKind = "codec"
)

// Error enriches transport errors with protocol, operation, and classification metadata.
type Error struct {
	Protocol Protocol
	Op       string
	Addr     string
	Kind     ErrorKind
	Err      error
}

// Error renders the enriched client error message.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Err == nil {
		return fmt.Sprintf("%s %s %s (%s)", e.Protocol, e.Op, e.Addr, e.Kind)
	}
	if e.Addr != "" {
		return fmt.Sprintf("%s %s %s (%s): %v", e.Protocol, e.Op, e.Addr, e.Kind, e.Err)
	}
	return fmt.Sprintf("%s %s (%s): %v", e.Protocol, e.Op, e.Kind, e.Err)
}

// Unwrap returns the underlying error.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// Timeout reports whether the error should be treated as a timeout.
func (e *Error) Timeout() bool {
	if e == nil {
		return false
	}
	if e.Kind == ErrorKindTimeout {
		return true
	}
	var netErr net.Error
	return errors.As(e.Err, &netErr) && netErr.Timeout()
}

// Temporary reports whether the error is marked as temporary.
func (e *Error) Temporary() bool {
	if e == nil {
		return false
	}
	if e.Kind == ErrorKindTemporary {
		return true
	}
	// net.Error.Temporary() 已弃用，这里仅检查 Kind 标记
	return false
}

// WrapError wraps err using the inferred ErrorKind.
func WrapError(protocol Protocol, op, addr string, err error) error {
	return WrapErrorWithKind(protocol, op, addr, "", err)
}

// WrapErrorWithKind wraps err with the supplied or inferred ErrorKind.
func WrapErrorWithKind(protocol Protocol, op, addr string, kind ErrorKind, err error) error {
	if err == nil {
		return nil
	}
	var existing *Error
	if errors.As(err, &existing) {
		return err
	}
	if protocol == "" {
		protocol = ProtocolUnknown
	}
	return &Error{
		Protocol: protocol,
		Op:       op,
		Addr:     addr,
		Kind:     lo.Ternary(kind != "", kind, classifyErrorKind(err)),
		Err:      err,
	}
}

// IsKind reports whether err is a *Error with the given kind.
func IsKind(err error, kind ErrorKind) bool {
	var e *Error
	if !errors.As(err, &e) {
		return false
	}
	return e.Kind == kind
}

// KindOf returns the ErrorKind carried by err when available.
func KindOf(err error) ErrorKind {
	var e *Error
	if !errors.As(err, &e) {
		return ErrorKindUnknown
	}
	return e.Kind
}

func classifyErrorKind(err error) ErrorKind {
	if err == nil {
		return ErrorKindUnknown
	}
	if kind, ok := classifyContextError(err); ok {
		return kind
	}
	if kind, ok := classifyClosedError(err); ok {
		return kind
	}
	if kind, ok := classifyTypedNetworkError(err); ok {
		return kind
	}
	return classifyMessageError(err)
}

func isConnRefused(err error) bool {
	if errors.Is(err, syscall.ECONNREFUSED) {
		return true
	}
	var errno syscall.Errno
	if errors.As(err, &errno) {
		return lo.Contains([]syscall.Errno{syscall.ECONNREFUSED, syscall.Errno(10061)}, errno)
	}
	var sysErr *os.SyscallError
	if errors.As(err, &sysErr) {
		return isConnRefused(sysErr.Err)
	}
	return false
}

func classifyContextError(err error) (ErrorKind, bool) {
	switch {
	case errors.Is(err, context.Canceled):
		return ErrorKindCanceled, true
	case errors.Is(err, context.DeadlineExceeded):
		return ErrorKindTimeout, true
	default:
		return "", false
	}
}

func classifyClosedError(err error) (ErrorKind, bool) {
	if errors.Is(err, net.ErrClosed) || errors.Is(err, os.ErrClosed) {
		return ErrorKindClosed, true
	}
	return "", false
}

func classifyTypedNetworkError(err error) (ErrorKind, bool) {
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return ErrorKindDNS, true
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return ErrorKindTimeout, true
	}

	if opErr, ok := errors.AsType[*net.OpError](err); ok {
		if opErr.Err != nil && isConnRefused(opErr.Err) {
			return ErrorKindConnRefused, true
		}
		return ErrorKindNetwork, true
	}

	if isConnRefused(err) {
		return ErrorKindConnRefused, true
	}

	if errors.As(err, &netErr) {
		return ErrorKindNetwork, true
	}

	return "", false
}

func classifyMessageError(err error) ErrorKind {
	msg := strings.ToLower(err.Error())
	switch {
	case strings.Contains(msg, "tls"), strings.Contains(msg, "x509"), strings.Contains(msg, "certificate"):
		return ErrorKindTLS
	case strings.Contains(msg, "use of closed network connection"), strings.Contains(msg, "file already closed"):
		return ErrorKindClosed
	case strings.Contains(msg, "network"):
		return ErrorKindNetwork
	default:
		return ErrorKindUnknown
	}
}
