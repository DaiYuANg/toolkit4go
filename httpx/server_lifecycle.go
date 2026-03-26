package httpx

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/DaiYuANg/arcgo/httpx/adapter"
)

// ListenPort starts related services on the provided port.
func (s *Server) ListenPort(port int) error {
	if port < 0 || port > 65535 {
		return fmt.Errorf("httpx: invalid port %d", port)
	}
	return s.Listen(fmt.Sprintf(":%d", port))
}

// Shutdown stops related services through the underlying adapter.
func (s *Server) Shutdown() error {
	if s == nil || s.adapter == nil {
		return fmt.Errorf("%w: adapter does not support shutdown", ErrAdapterNotFound)
	}
	if s.logger != nil && s.logger.Enabled(context.Background(), slog.LevelDebug) {
		s.logger.Debug("httpx shutdown requested",
			slog.String("adapter", s.adapter.Name()),
		)
	}

	var shutdownErr error
	if useHostCapability(s, func(shutdownable adapter.ShutdownAdapter) {
		shutdownErr = shutdownable.Shutdown()
	}) {
		if shutdownErr != nil && s.logger != nil {
			s.logger.Error("httpx shutdown failed",
				slog.String("adapter", s.adapter.Name()),
				slog.String("error", shutdownErr.Error()),
			)
		} else if s.logger != nil && s.logger.Enabled(context.Background(), slog.LevelDebug) {
			s.logger.Debug("httpx shutdown completed",
				slog.String("adapter", s.adapter.Name()),
			)
		}
		return shutdownErr
	}

	return fmt.Errorf("%w: adapter %q does not support shutdown", ErrAdapterNotFound, s.adapter.Name())
}

// Listen starts related services.
func (s *Server) Listen(addr string) error {
	s.freezeConfiguration()

	name := "unknown"
	if s.adapter != nil {
		name = s.adapter.Name()
	}

	s.logger.Info("Starting server",
		slog.String("address", addr),
		slog.String("adapter", name),
		slog.Int("routes", s.RouteCount()),
	)
	if s.logger.Enabled(context.Background(), slog.LevelDebug) {
		s.logger.Debug("httpx listen entering adapter",
			slog.String("address", addr),
			slog.String("adapter", name),
		)
	}

	var listenErr error
	if useHostCapability(s, func(listenable adapter.ListenableAdapter) {
		listenErr = listenable.Listen(addr)
	}) {
		if listenErr != nil {
			s.logger.Error("httpx listen failed",
				slog.String("address", addr),
				slog.String("adapter", name),
				slog.String("error", listenErr.Error()),
			)
			return fmt.Errorf("httpx: adapter %q listen on %q: %w", name, addr, listenErr)
		}
		if s.logger.Enabled(context.Background(), slog.LevelDebug) {
			s.logger.Debug("httpx listen returned",
				slog.String("address", addr),
				slog.String("adapter", name),
			)
		}
		return nil
	}

	return fmt.Errorf("%w: adapter %q does not support direct listening", ErrAdapterNotFound, name)
}

// ListenAndServe starts related services.
func (s *Server) ListenAndServe(addr string) error {
	return s.Listen(addr)
}

// ListenAndServeContext starts related services.
func (s *Server) ListenAndServeContext(ctx context.Context, addr string) error {
	s.freezeConfiguration()

	name := "unknown"
	if s != nil && s.adapter != nil {
		name = s.adapter.Name()
	}

	var listenErr error
	if useHostCapability(s, func(listenable adapter.ContextListenableAdapter) {
		listenErr = listenable.ListenContext(ctx, addr)
	}) {
		return listenErr
	}

	var listenable adapter.ListenableAdapter
	var shutdownable adapter.ShutdownAdapter
	if !useHostCapability(s, func(host adapter.ListenableAdapter) {
		listenable = host
	}) || !useHostCapability(s, func(host adapter.ShutdownAdapter) {
		shutdownable = host
	}) {
		return fmt.Errorf("%w: adapter %q does not support context listening", ErrAdapterNotFound, name)
	}

	if ctx == nil {
		ctx = context.Background()
	}

	s.logger.Info("Starting server with context",
		slog.String("address", addr),
		slog.String("adapter", name),
	)
	if s.logger.Enabled(context.Background(), slog.LevelDebug) {
		s.logger.Debug("httpx context listen using fallback adapter path",
			slog.String("address", addr),
			slog.String("adapter", name),
		)
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- listenable.Listen(addr)
	}()

	select {
	case err := <-errCh:
		if err != nil && s.logger != nil {
			s.logger.Error("httpx context listen exited with error",
				slog.String("address", addr),
				slog.String("adapter", name),
				slog.String("error", err.Error()),
			)
		} else if s.logger.Enabled(context.Background(), slog.LevelDebug) {
			s.logger.Debug("httpx context listen exited",
				slog.String("address", addr),
				slog.String("adapter", name),
			)
		}
		return err
	case <-ctx.Done():
		if s.logger.Enabled(context.Background(), slog.LevelDebug) {
			s.logger.Debug("httpx context cancelled; shutting down adapter",
				slog.String("address", addr),
				slog.String("adapter", name),
			)
		}
		if err := shutdownable.Shutdown(); err != nil && !errors.Is(err, ErrAdapterNotFound) {
			s.logger.Error("httpx context shutdown failed",
				slog.String("address", addr),
				slog.String("adapter", name),
				slog.String("error", err.Error()),
			)
			return fmt.Errorf("httpx: shutdown server on %q: %w", addr, err)
		}
		select {
		case err := <-errCh:
			if err != nil && s.logger != nil {
				s.logger.Error("httpx listen returned after context cancellation",
					slog.String("address", addr),
					slog.String("adapter", name),
					slog.String("error", err.Error()),
				)
			} else if s.logger.Enabled(context.Background(), slog.LevelDebug) {
				s.logger.Debug("httpx listen returned after context cancellation",
					slog.String("address", addr),
					slog.String("adapter", name),
				)
			}
			return err
		default:
			if s.logger.Enabled(context.Background(), slog.LevelDebug) {
				s.logger.Debug("httpx context fallback returning context error",
					slog.String("address", addr),
					slog.String("adapter", name),
				)
			}
			return ctx.Err()
		}
	}
}
