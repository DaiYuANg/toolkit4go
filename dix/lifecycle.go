package dix

import (
	"context"
	"fmt"
	"log/slog"
)

// StartHook is executed when the application starts.
type StartHook func(ctx context.Context) error

// StopHook is executed when the application stops.
type StopHook func(ctx context.Context) error

// Lifecycle manages application lifecycle hooks.
type Lifecycle interface {
	OnStart(hook StartHook)
	OnStop(hook StopHook)
}

// lifecycleImpl is the internal implementation.
type lifecycleImpl struct {
	startHooks []StartHook
	stopHooks  []StopHook
	logger     *slog.Logger
}

func newLifecycle() *lifecycleImpl {
	return &lifecycleImpl{
		startHooks: make([]StartHook, 0),
		stopHooks:  make([]StopHook, 0),
		logger:     slog.Default(),
	}
}

func (l *lifecycleImpl) OnStart(hook StartHook) {
	l.startHooks = append(l.startHooks, hook)
}

func (l *lifecycleImpl) OnStop(hook StopHook) {
	l.stopHooks = append(l.stopHooks, hook)
}

func (l *lifecycleImpl) executeStartHooks(ctx context.Context, _ *Container) error {
	for i, hook := range l.startHooks {
		if err := hook(ctx); err != nil {
			if l.logger != nil {
				l.logger.Error("start hook failed", "index", i, "error", err)
			}
			return fmt.Errorf("start hook %d failed: %w", i, err)
		}
	}
	return nil
}

func (l *lifecycleImpl) executeStopHooks(ctx context.Context, _ *Container) error {
	for i := len(l.stopHooks) - 1; i >= 0; i-- {
		hook := l.stopHooks[i]
		if err := hook(ctx); err != nil {
			if l.logger != nil {
				l.logger.Error("stop hook failed", "index", i, "error", err)
			}
			return fmt.Errorf("stop hook %d failed: %w", i, err)
		}
	}
	return nil
}

// OnStartHook creates a typed start hook.
// The container is captured in the closure at registration time.
func OnStartHook[T any](c *Container, fn func(context.Context, T) error) func(lc Lifecycle) {
	return func(lc Lifecycle) {
		lc.OnStart(func(ctx context.Context) error {
			t, err := ResolveAs[T](c)
			if err != nil {
				return fmt.Errorf("resolving dependency: %w", err)
			}
			return fn(ctx, t)
		})
	}
}

// OnStopHook creates a typed stop hook.
func OnStopHook[T any](c *Container, fn func(context.Context, T) error) func(lc Lifecycle) {
	return func(lc Lifecycle) {
		lc.OnStop(func(ctx context.Context) error {
			t, err := ResolveAs[T](c)
			if err != nil {
				return fmt.Errorf("resolving dependency: %w", err)
			}
			return fn(ctx, t)
		})
	}
}

// OnStartHook2 creates a typed start hook with 2 dependencies.
func OnStartHook2[T1, T2 any](c *Container, fn func(context.Context, T1, T2) error) func(lc Lifecycle) {
	return func(lc Lifecycle) {
		lc.OnStart(func(ctx context.Context) error {
			t1, err := ResolveAs[T1](c)
			if err != nil {
				return err
			}
			t2, err := ResolveAs[T2](c)
			if err != nil {
				return err
			}
			return fn(ctx, t1, t2)
		})
	}
}

// OnStopHook2 creates a typed stop hook with 2 dependencies.
func OnStopHook2[T1, T2 any](c *Container, fn func(context.Context, T1, T2) error) func(lc Lifecycle) {
	return func(lc Lifecycle) {
		lc.OnStop(func(ctx context.Context) error {
			t1, err := ResolveAs[T1](c)
			if err != nil {
				return err
			}
			t2, err := ResolveAs[T2](c)
			if err != nil {
				return err
			}
			return fn(ctx, t1, t2)
		})
	}
}

// OnStartHook3 creates a typed start hook with 3 dependencies.
func OnStartHook3[T1, T2, T3 any](c *Container, fn func(context.Context, T1, T2, T3) error) func(lc Lifecycle) {
	return func(lc Lifecycle) {
		lc.OnStart(func(ctx context.Context) error {
			t1, err := ResolveAs[T1](c)
			if err != nil {
				return err
			}
			t2, err := ResolveAs[T2](c)
			if err != nil {
				return err
			}
			t3, err := ResolveAs[T3](c)
			if err != nil {
				return err
			}
			return fn(ctx, t1, t2, t3)
		})
	}
}

// OnStopHook3 creates a typed stop hook with 3 dependencies.
func OnStopHook3[T1, T2, T3 any](c *Container, fn func(context.Context, T1, T2, T3) error) func(lc Lifecycle) {
	return func(lc Lifecycle) {
		lc.OnStop(func(ctx context.Context) error {
			t1, err := ResolveAs[T1](c)
			if err != nil {
				return err
			}
			t2, err := ResolveAs[T2](c)
			if err != nil {
				return err
			}
			t3, err := ResolveAs[T3](c)
			if err != nil {
				return err
			}
			return fn(ctx, t1, t2, t3)
		})
	}
}

// Hook is deprecated.
var Hook = hookHelper{}

type hookHelper struct{}

func (hookHelper) OnStart(any) {
	panic("use OnStartHook[T](c, fn)(lc) instead")
}

func (hookHelper) OnStop(any) {
	panic("use OnStopHook[T](c, fn)(lc) instead")
}
