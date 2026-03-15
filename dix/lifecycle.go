package dix

import (
	"context"
	"fmt"
)

// Lifecycle manages application lifecycle hooks.
type Lifecycle interface {
	OnStart(hook any)
	OnStop(hook any)
}

// lifecycleImpl is the internal implementation.
type lifecycleImpl struct {
	startHooks []func(ctx context.Context) error
	stopHooks  []func(ctx context.Context) error
}

func newLifecycle() *lifecycleImpl {
	return &lifecycleImpl{
		startHooks: make([]func(ctx context.Context) error, 0),
		stopHooks:  make([]func(ctx context.Context) error, 0),
	}
}

func (l *lifecycleImpl) OnStart(hook any) {
	// The hook is expected to be func(ctx context.Context) error
	if fn, ok := hook.(func(context.Context) error); ok {
		l.startHooks = append(l.startHooks, fn)
	}
}

func (l *lifecycleImpl) OnStop(hook any) {
	if fn, ok := hook.(func(context.Context) error); ok {
		l.stopHooks = append(l.stopHooks, fn)
	}
}

func (l *lifecycleImpl) executeStartHooks(ctx context.Context, _ *Container) error {
	for i, hook := range l.startHooks {
		if err := hook(ctx); err != nil {
			return fmt.Errorf("start hook %d failed: %w", i, err)
		}
	}
	return nil
}

func (l *lifecycleImpl) executeStopHooks(ctx context.Context, _ *Container) error {
	for i := len(l.stopHooks) - 1; i >= 0; i-- {
		hook := l.stopHooks[i]
		if err := hook(ctx); err != nil {
			return fmt.Errorf("stop hook %d failed: %w", i, err)
		}
	}
	return nil
}

// OnStartHook creates a typed start hook.
// The container is captured in the closure at registration time.
func OnStartHook[T any](c *Container, fn func(T) error) func(lc Lifecycle) {
	return func(lc Lifecycle) {
		lc.OnStart(func(ctx context.Context) error {
			t, err := ResolveAs[T](c)
			if err != nil {
				return fmt.Errorf("resolving dependency: %w", err)
			}
			return fn(t)
		})
	}
}

// OnStopHook creates a typed stop hook.
func OnStopHook[T any](c *Container, fn func(T) error) func(lc Lifecycle) {
	return func(lc Lifecycle) {
		lc.OnStop(func(ctx context.Context) error {
			t, err := ResolveAs[T](c)
			if err != nil {
				return fmt.Errorf("resolving dependency: %w", err)
			}
			return fn(t)
		})
	}
}

// OnStartHook2 creates a typed start hook with 2 dependencies.
func OnStartHook2[T1, T2 any](c *Container, fn func(T1, T2) error) func(lc Lifecycle) {
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
			return fn(t1, t2)
		})
	}
}

// OnStopHook2 creates a typed stop hook with 2 dependencies.
func OnStopHook2[T1, T2 any](c *Container, fn func(T1, T2) error) func(lc Lifecycle) {
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
			return fn(t1, t2)
		})
	}
}

// OnStartHook3 creates a typed start hook with 3 dependencies.
func OnStartHook3[T1, T2, T3 any](c *Container, fn func(T1, T2, T3) error) func(lc Lifecycle) {
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
			return fn(t1, t2, t3)
		})
	}
}

// OnStopHook3 creates a typed stop hook with 3 dependencies.
func OnStopHook3[T1, T2, T3 any](c *Container, fn func(T1, T2, T3) error) func(lc Lifecycle) {
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
			return fn(t1, t2, t3)
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
