package std

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/DaiYuANg/arcgo/httpx/adapter"
	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"
)

// Adapter implements the std/chi runtime bridge for httpx.
type Adapter struct {
	router *chi.Mux
	prefix string
	logger *slog.Logger
	huma   huma.API
	docs   *adapter.DocsController
	server ServerOptions
}

// New constructs a std adapter backed by a chi router and Huma API.
func New(opts ...adapter.HumaOptions) *Adapter {
	options := DefaultOptions()
	options.Huma = adapter.MergeHumaOptions(opts...)
	return NewWithOptions(options)
}

// WithLogger replaces the adapter logger.
func (a *Adapter) WithLogger(logger *slog.Logger) *Adapter {
	a.SetLogger(logger)
	return a
}

// SetLogger replaces the adapter logger.
func (a *Adapter) SetLogger(logger *slog.Logger) {
	if a == nil || logger == nil {
		return
	}
	a.logger = logger
}

// Name returns the adapter name.
func (a *Adapter) Name() string {
	return "std"
}

// Handle registers a native handler on the chi router.
func (a *Adapter) Handle(method, path string, handler adapter.HandlerFunc) {
	fullPath := joinPath(a.prefix, path)
	a.router.Method(method, fullPath, a.wrapHandler(handler))
}

// Group returns a prefixed child adapter that shares the same router and Huma API.
func (a *Adapter) Group(prefix string) adapter.Adapter {
	nextPrefix := a.prefix
	if prefix != "" && prefix != "/" {
		nextPrefix = joinPath(a.prefix, prefix)
	}
	return &Adapter{
		router: a.router,
		prefix: nextPrefix,
		logger: a.logger,
		huma:   a.huma,
		docs:   a.docs,
		server: a.server,
	}
}

// ServeHTTP serves docs routes first, then falls through to the chi router.
func (a *Adapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if a.docs != nil && a.docs.ServeHTTP(w, r) {
		return
	}
	a.router.ServeHTTP(w, r)
}

// Router exposes the underlying chi router.
func (a *Adapter) Router() *chi.Mux {
	return a.router
}

// Listen starts related services.
func (a *Adapter) Listen(addr string) error {
	if err := a.httpServer(addr).ListenAndServe(); err != nil {
		return fmt.Errorf("httpx/std: listen on %q: %w", addr, err)
	}
	return nil
}

// ListenContext starts related services.
func (a *Adapter) ListenContext(ctx context.Context, addr string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	server := a.httpServer(addr)

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("httpx/std: listen on %q: %w", addr, err)
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), a.server.ShutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("httpx/std: shutdown on %q: %w", addr, err)
		}
		err := <-errCh
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("httpx/std: listen on %q: %w", addr, err)
	}
}

// wrapHandler converts adapter-native handlers into `http.HandlerFunc`.
func (a *Adapter) wrapHandler(handler adapter.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(r.Context(), w, r); err != nil {
			a.logger.Error("Handler error",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("error", err.Error()),
			)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// HumaAPI exposes the underlying Huma API.
func (a *Adapter) HumaAPI() huma.API {
	return a.huma
}

// ConfigureHumaOptions updates adapter-managed docs/openapi routing.
func (a *Adapter) ConfigureHumaOptions(opts adapter.HumaOptions) {
	if a == nil || a.docs == nil {
		return
	}
	a.docs.Configure(opts)
}

func (a *Adapter) httpServer(addr string) *http.Server {
	return &http.Server{
		Addr:           addr,
		Handler:        a,
		ReadTimeout:    a.server.ReadTimeout,
		WriteTimeout:   a.server.WriteTimeout,
		IdleTimeout:    a.server.IdleTimeout,
		MaxHeaderBytes: a.server.MaxHeaderBytes,
	}
}
