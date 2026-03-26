package kvx

import "log/slog"

func LogDebug(logger *slog.Logger, debug bool, msg string, attrs ...any) {
	if logger == nil || !debug {
		return
	}
	logger.Debug(msg, attrs...)
}

func LogError(logger *slog.Logger, msg string, attrs ...any) {
	if logger == nil {
		return
	}
	logger.Error(msg, attrs...)
}
