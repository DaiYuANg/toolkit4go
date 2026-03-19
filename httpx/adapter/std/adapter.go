package std

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"
)

type lifecycleState struct {
	mu     sync.Mutex
	server *http.Server
}

// Adapter implements the std/chi runtime bridge for httpx.
type Adapter struct {
	router    *chi.Mux
	huma      huma.API
	lifecycle *lifecycleState
}

// Name returns the adapter name.
func (a *Adapter) Name() string {
	return "std"
}

// Router exposes the underlying chi router.
func (a *Adapter) Router() *chi.Mux {
	return a.router
}

// Listen starts related services.
func (a *Adapter) Listen(addr string) error {
	server := a.httpServer(addr)
	release := a.trackServer(server)
	defer release()

	if err := server.ListenAndServe(); err != nil {
		return fmt.Errorf("httpx/std: listen on %q: %w", addr, err)
	}
	return nil
}

// Shutdown stops the active std server.
func (a *Adapter) Shutdown() error {
	server := a.activeServer()
	if server == nil {
		return nil
	}
	if err := server.Shutdown(context.Background()); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("httpx/std: shutdown: %w", err)
	}
	return nil
}

// ListenContext starts related services.
func (a *Adapter) ListenContext(ctx context.Context, addr string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	server := a.httpServer(addr)
	release := a.trackServer(server)
	defer release()

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
		if err := a.Shutdown(); err != nil {
			return fmt.Errorf("httpx/std: shutdown on %q: %w", addr, err)
		}
		err := <-errCh
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("httpx/std: listen on %q: %w", addr, err)
	}
}

// HumaAPI exposes the underlying Huma API.
func (a *Adapter) HumaAPI() huma.API {
	return a.huma
}

func (a *Adapter) httpServer(addr string) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: a.router,
	}
}

func (a *Adapter) trackServer(server *http.Server) func() {
	if a == nil || a.lifecycle == nil {
		return func() {}
	}

	a.lifecycle.mu.Lock()
	a.lifecycle.server = server
	a.lifecycle.mu.Unlock()

	return func() {
		a.lifecycle.mu.Lock()
		if a.lifecycle.server == server {
			a.lifecycle.server = nil
		}
		a.lifecycle.mu.Unlock()
	}
}

func (a *Adapter) activeServer() *http.Server {
	if a == nil || a.lifecycle == nil {
		return nil
	}

	a.lifecycle.mu.Lock()
	defer a.lifecycle.mu.Unlock()
	return a.lifecycle.server
}
