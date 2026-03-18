package std

import (
	"log/slog"
	"time"

	"github.com/DaiYuANg/arcgo/httpx/adapter"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

// ServerOptions configures the std adapter's underlying http.Server.
type ServerOptions struct {
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	MaxHeaderBytes  int
}

// DefaultServerOptions returns the default std adapter server config.
func DefaultServerOptions() ServerOptions {
	return ServerOptions{
		ReadTimeout:     15 * time.Second,
		WriteTimeout:    15 * time.Second,
		IdleTimeout:     60 * time.Second,
		ShutdownTimeout: 5 * time.Second,
		MaxHeaderBytes:  1 << 20,
	}
}

// Options configures std adapter construction.
type Options struct {
	Huma   adapter.HumaOptions
	Logger *slog.Logger
	Server ServerOptions
}

// DefaultOptions returns the default std adapter config.
func DefaultOptions() Options {
	return Options{
		Huma:   adapter.DefaultHumaOptions(),
		Logger: slog.Default(),
		Server: DefaultServerOptions(),
	}
}

// NewWithOptions constructs a std adapter from explicit construction-time options.
func NewWithOptions(opts Options) *Adapter {
	router := chi.NewMux()

	humaOpts := adapter.MergeHumaOptions(opts.Huma)
	cfg := huma.DefaultConfig(humaOpts.Title, humaOpts.Version)
	adapter.ApplyHumaConfig(&cfg, humaOpts)

	docsCfg := cfg
	docsCfg.DocsPath = ""
	docsCfg.OpenAPIPath = ""
	docsCfg.SchemasPath = ""
	api := humachi.New(router, docsCfg)
	docs := adapter.NewDocsController(api, humaOpts)

	return &Adapter{
		router: router,
		logger: defaultLogger(opts.Logger),
		huma:   api,
		docs:   docs,
		server: mergeServerOptions(opts.Server),
	}
}
