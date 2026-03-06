package httpx

import (
	"log/slog"

	"github.com/DaiYuANg/arcgo/httpx/adapter"
	"github.com/go-playground/validator/v10"
)

// WithAdapter configures related behavior.
func WithAdapter(adapter adapter.Adapter) ServerOption {
	return func(s *Server) {
		s.adapter = adapter
	}
}

// WithBasePath configures related behavior.
func WithBasePath(path string) ServerOption {
	return func(s *Server) {
		s.basePath = normalizeRoutePrefix(path)
	}
}

// WithLogger configures related behavior.
func WithLogger(logger *slog.Logger) ServerOption {
	return func(s *Server) {
		s.logger = logger
	}
}

// WithPrintRoutes configures related behavior.
func WithPrintRoutes(enabled bool) ServerOption {
	return func(s *Server) {
		s.printRoutes = enabled
	}
}

// WithOpenAPIInfo configures related behavior.
// Note: This option must be passed to the adapter's New function directly.
// This function is kept for backward compatibility but has no effect when used with NewServer.
// Use adapter-level options instead: adapter.New(engine, adapter.HumaOptions{Title: "...", Version: "..."})
func WithOpenAPIInfo(title, version, description string) ServerOption {
	return func(s *Server) {
		// No-op: Huma options are now configured at adapter creation time
		_ = title
		_ = version
		_ = description
	}
}

// WithOpenAPIDocs provides default behavior.
// Note: This option must be passed to the adapter's New function directly.
// This function is kept for backward compatibility but has no effect when used with NewServer.
// Use adapter-level options instead: adapter.New(engine, adapter.HumaOptions{DisableDocsRoutes: true})
func WithOpenAPIDocs(enabled bool) ServerOption {
	return func(s *Server) {
		// No-op: Huma options are now configured at adapter creation time
		_ = enabled
	}
}

// WithValidation enables related functionality.
func WithValidation() ServerOption {
	return func(s *Server) {
		if s.validator == nil {
			s.validator = validator.New(validator.WithRequiredStructEnabled())
		}
	}
}

// WithValidator closes related resources.
func WithValidator(v *validator.Validate) ServerOption {
	return func(s *Server) {
		s.validator = v
	}
}

