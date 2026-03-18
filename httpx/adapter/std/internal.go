package std

import (
	"log/slog"
	"strings"
)

func defaultLogger(logger *slog.Logger) *slog.Logger {
	if logger != nil {
		return logger
	}
	return slog.Default()
}

func mergeServerOptions(opts ServerOptions) ServerOptions {
	defaults := DefaultServerOptions()
	if opts.ReadTimeout > 0 {
		defaults.ReadTimeout = opts.ReadTimeout
	}
	if opts.WriteTimeout > 0 {
		defaults.WriteTimeout = opts.WriteTimeout
	}
	if opts.IdleTimeout > 0 {
		defaults.IdleTimeout = opts.IdleTimeout
	}
	if opts.ShutdownTimeout > 0 {
		defaults.ShutdownTimeout = opts.ShutdownTimeout
	}
	if opts.MaxHeaderBytes > 0 {
		defaults.MaxHeaderBytes = opts.MaxHeaderBytes
	}
	return defaults
}

func joinPath(prefix, path string) string {
	cleanPrefix := strings.TrimRight(prefix, "/")
	if cleanPrefix == "" {
		if path == "" {
			return "/"
		}
		if strings.HasPrefix(path, "/") {
			return path
		}
		return "/" + path
	}

	if path == "" || path == "/" {
		return cleanPrefix
	}
	if strings.HasPrefix(path, "/") {
		return cleanPrefix + path
	}
	return cleanPrefix + "/" + path
}
