package httpx

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	collectionlist "github.com/DaiYuANg/arcgo/collectionx/list"
	collectionset "github.com/DaiYuANg/arcgo/collectionx/set"
	"github.com/DaiYuANg/arcgo/httpx/adapter"
	"github.com/DaiYuANg/arcgo/httpx/adapter/std"
	"github.com/danielgtaylor/huma/v2"
	"github.com/go-playground/validator/v10"
	"github.com/samber/lo"
)

// Server documents related behavior.
//
// Note.
// Note.
// Note.
// Note.
// Note.
type Server struct {
	adapter     adapter.Adapter
	basePath    string
	routesMu    sync.RWMutex
	routes      *collectionlist.List[RouteInfo]
	routeKeys   *collectionset.Set[string]
	logger      *slog.Logger
	printRoutes bool
	validator   *validator.Validate
}

// Group documents related behavior.
type Group struct {
	server *Server
	prefix string
}

// ServerOption documents related behavior.
type ServerOption func(*Server)

// NewServer creates related functionality.
func NewServer(opts ...ServerOption) *Server {
	s := &Server{
		logger:    slog.Default(),
		routes:    collectionlist.NewList[RouteInfo](),
		routeKeys: collectionset.NewSet[string](),
	}

	lo.ForEach(opts, func(opt ServerOption, _ int) {
		opt(s)
	})

	if s.adapter == nil {
		s.adapter = std.New()
	}

	return s
}

// Group creates related functionality.
func (s *Server) Group(prefix string) *Group {
	return &Group{
		server: s,
		prefix: normalizeRoutePrefix(prefix),
	}
}

// printRoutesIfEnabled documents related behavior.
func (s *Server) printRoutesIfEnabled() {
	if !s.printRoutes {
		return
	}

	routes := s.routesSnapshot()
	s.logger.Info("Registered routes", slog.Int("count", len(routes)))
	lo.ForEach(routes, func(route RouteInfo, _ int) {
		s.logger.Info("  "+route.String(),
			slog.String("method", route.Method),
			slog.String("path", route.Path),
			slog.String("handler", route.HandlerName),
		)
	})
}

// GetRoutes returns related data.
func (s *Server) GetRoutes() []RouteInfo {
	return s.routesSnapshot()
}

// GetRoutesByMethod documents related behavior.
func (s *Server) GetRoutesByMethod(method string) []RouteInfo {
	method = strings.ToUpper(method)
	return lo.Filter(s.routesSnapshot(), func(route RouteInfo, _ int) bool {
		return route.Method == method
	})
}

// GetRoutesByPath documents related behavior.
func (s *Server) GetRoutesByPath(prefix string) []RouteInfo {
	if prefix == "" {
		return s.routesSnapshot()
	}
	return lo.Filter(s.routesSnapshot(), func(route RouteInfo, _ int) bool {
		return strings.HasPrefix(route.Path, prefix)
	})
}

// HasRoute checks related state.
func (s *Server) HasRoute(method, path string) bool {
	s.routesMu.RLock()
	defer s.routesMu.RUnlock()
	return s.routeKeys.Contains(routeKey(strings.ToUpper(method), path))
}

// RouteCount returns related data.
func (s *Server) RouteCount() int {
	s.routesMu.RLock()
	defer s.routesMu.RUnlock()
	return s.routes.Len()
}

// Handler returns related data.
func (s *Server) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.adapter.ServeHTTP(w, r)
	})
}

// ServeHTTP documents related behavior.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Handler().ServeHTTP(w, r)
}

// ListenAndServe starts related services.
func (s *Server) ListenAndServe(addr string) error {
	routeCount := s.RouteCount()
	s.logger.Info("Starting server",
		slog.String("address", addr),
		slog.String("adapter", s.adapter.Name()),
		slog.Int("routes", routeCount),
	)

	if listenable, ok := s.adapter.(adapter.ListenableAdapter); ok {
		if err := listenable.Listen(addr); err != nil {
			return fmt.Errorf("httpx: adapter %q listen on %q: %w", s.adapter.Name(), addr, err)
		}
		return nil
	}

	if err := http.ListenAndServe(addr, s.Handler()); err != nil {
		return fmt.Errorf("httpx: listen on %q: %w", addr, err)
	}
	return nil
}

// ListenAndServeContext starts related services.
func (s *Server) ListenAndServeContext(ctx context.Context, addr string) error {
	if listenable, ok := s.adapter.(adapter.ContextListenableAdapter); ok {
		return listenable.ListenContext(ctx, addr)
	}

	server := &http.Server{
		Addr:    addr,
		Handler: s.Handler(),
	}

	s.logger.Info("Starting server with context", slog.String("address", addr))

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("httpx: listen on %q: %w", addr, err)
	case <-ctx.Done():
		s.logger.Info("Shutting down server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("httpx: shutdown server on %q: %w", addr, err)
		}
		err := <-errCh
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("httpx: listen on %q: %w", addr, err)
	}
}

// Adapter returns related data.
func (s *Server) Adapter() adapter.Adapter {
	return s.adapter
}

// Logger returns related data.
func (s *Server) Logger() *slog.Logger {
	return s.logger
}

// Validator returns related data.
func (s *Server) Validator() *validator.Validate {
	return s.validator
}

// HumaAPI returns related data.
func (s *Server) HumaAPI() huma.API {
	return s.adapter.HumaAPI()
}

func (s *Server) addRoute(route RouteInfo) {
	s.routesMu.Lock()
	defer s.routesMu.Unlock()

	key := routeKey(route.Method, route.Path)
	if s.routeKeys.Contains(key) {
		return
	}

	s.routeKeys.Add(key)
	s.routes.Add(route)
}

func (s *Server) routesSnapshot() []RouteInfo {
	s.routesMu.RLock()
	defer s.routesMu.RUnlock()
	return s.routes.Values()
}

func routeKey(method, path string) string {
	return method + " " + path
}
