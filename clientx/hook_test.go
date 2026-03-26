package clientx

import (
	"context"
	"log/slog"
	"testing"
)

func TestHookFuncsDispatch(t *testing.T) {
	var dialCalled bool
	var ioCalled bool

	h := HookFuncs{
		OnDialFunc: func(event DialEvent) {
			dialCalled = event.Protocol == ProtocolTCP
		},
		OnIOFunc: func(event IOEvent) {
			ioCalled = event.Protocol == ProtocolHTTP
		},
	}

	EmitDial([]Hook{h}, DialEvent{Protocol: ProtocolTCP})
	EmitIO([]Hook{h}, IOEvent{Protocol: ProtocolHTTP})

	if !dialCalled {
		t.Fatal("expected dial hook to be called")
	}
	if !ioCalled {
		t.Fatal("expected io hook to be called")
	}
}

func TestEmitHookPanicIsolation(t *testing.T) {
	dialCalled := false
	ioCalled := false

	hooks := []Hook{
		HookFuncs{
			OnDialFunc: func(event DialEvent) {
				panic("dial hook panic")
			},
			OnIOFunc: func(event IOEvent) {
				panic("io hook panic")
			},
		},
		HookFuncs{
			OnDialFunc: func(event DialEvent) {
				dialCalled = true
			},
			OnIOFunc: func(event IOEvent) {
				ioCalled = true
			},
		},
	}

	EmitDial(hooks, DialEvent{Protocol: ProtocolTCP})
	EmitIO(hooks, IOEvent{Protocol: ProtocolHTTP})

	if !dialCalled {
		t.Fatal("expected subsequent dial hook to be called after panic")
	}
	if !ioCalled {
		t.Fatal("expected subsequent io hook to be called after panic")
	}
}

type memoryLogHandler struct {
	records []memoryLogRecord
}

type memoryLogRecord struct {
	level   slog.Level
	message string
	attrs   map[string]any
}

func (h *memoryLogHandler) Enabled(context.Context, slog.Level) bool { return true }

func (h *memoryLogHandler) Handle(_ context.Context, record slog.Record) error {
	entry := memoryLogRecord{
		level:   record.Level,
		message: record.Message,
		attrs:   map[string]any{},
	}
	record.Attrs(func(attr slog.Attr) bool {
		entry.attrs[attr.Key] = attr.Value.Any()
		return true
	})
	h.records = append(h.records, entry)
	return nil
}

func (h *memoryLogHandler) WithAttrs([]slog.Attr) slog.Handler { return h }
func (h *memoryLogHandler) WithGroup(string) slog.Handler      { return h }

func TestLoggingHookEmitsDialAndIORecords(t *testing.T) {
	handler := &memoryLogHandler{}
	logger := slog.New(handler)
	hook := NewLoggingHook(logger)

	EmitDial([]Hook{hook}, DialEvent{Protocol: ProtocolTCP, Op: "dial", Network: "tcp", Addr: "127.0.0.1:9000"})
	EmitIO([]Hook{hook}, IOEvent{Protocol: ProtocolHTTP, Op: "get", Addr: "http://example.com", Bytes: 32})

	if len(handler.records) != 2 {
		t.Fatalf("expected 2 log records, got %d", len(handler.records))
	}
	if handler.records[0].message != "clientx dial" {
		t.Fatalf("unexpected dial log message: %s", handler.records[0].message)
	}
	if handler.records[1].message != "clientx io" {
		t.Fatalf("unexpected io log message: %s", handler.records[1].message)
	}
}
