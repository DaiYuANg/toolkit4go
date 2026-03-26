package clientx

import "log/slog"

type LoggingHookOption func(*loggingHookConfig)

type loggingHookConfig struct {
	includeAddress bool
}

func WithLoggingHookAddress(enabled bool) LoggingHookOption {
	return func(cfg *loggingHookConfig) {
		cfg.includeAddress = enabled
	}
}

func NewLoggingHook(logger *slog.Logger, opts ...LoggingHookOption) Hook {
	cfg := loggingHookConfig{includeAddress: true}
	Apply(&cfg, opts...)
	if logger == nil {
		logger = slog.Default()
	}
	return &loggingHook{logger: logger, cfg: cfg}
}

type loggingHook struct {
	logger *slog.Logger
	cfg    loggingHookConfig
}

func (h *loggingHook) OnDial(event DialEvent) {
	if h == nil || h.logger == nil {
		return
	}
	attrs := []any{
		"protocol", event.Protocol,
		"op", event.Op,
		"network", event.Network,
		"duration", event.Duration,
	}
	if h.cfg.includeAddress && event.Addr != "" {
		attrs = append(attrs, "addr", event.Addr)
	}
	if event.Err != nil {
		attrs = append(attrs, "error", event.Err, "error_kind", KindOf(event.Err))
		h.logger.Error("clientx dial", attrs...)
		return
	}
	h.logger.Debug("clientx dial", attrs...)
}

func (h *loggingHook) OnIO(event IOEvent) {
	if h == nil || h.logger == nil {
		return
	}
	attrs := []any{
		"protocol", event.Protocol,
		"op", event.Op,
		"bytes", event.Bytes,
		"duration", event.Duration,
	}
	if h.cfg.includeAddress && event.Addr != "" {
		attrs = append(attrs, "addr", event.Addr)
	}
	if event.Err != nil {
		attrs = append(attrs, "error", event.Err, "error_kind", KindOf(event.Err))
		h.logger.Error("clientx io", attrs...)
		return
	}
	h.logger.Debug("clientx io", attrs...)
}
