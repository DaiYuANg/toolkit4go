package dix

import (
	"context"
	"fmt"

	do "github.com/samber/do/v2"
)

// Container wraps samber/do.Injector.
type Container struct {
	injector do.Injector
}

func newContainer() *Container {
	return &Container{
		injector: do.New(),
	}
}

// Injector returns the underlying do.Injector.
func (c *Container) Injector() do.Injector {
	return c.injector
}

// Shutdown shuts down the container.
func (c *Container) Shutdown(ctx context.Context) error {
	c.injector.Shutdown()
	return nil
}

// ProvideT registers a typed provider with no dependencies.
func ProvideT[T any](c *Container, fn func() T) {
	do.Provide(c.injector, func(i do.Injector) (T, error) {
		return fn(), nil
	})
}

// Provide1T registers a typed provider with 1 dependency.
func Provide1T[T, D1 any](c *Container, fn func(D1) T) {
	do.Provide(c.injector, func(i do.Injector) (T, error) {
		d1, err := do.Invoke[D1](i)
		if err != nil {
			var zero T
			return zero, err
		}
		return fn(d1), nil
	})
}

// Provide2T registers a typed provider with 2 dependencies.
func Provide2T[T, D1, D2 any](c *Container, fn func(D1, D2) T) {
	do.Provide(c.injector, func(i do.Injector) (T, error) {
		d1, err := do.Invoke[D1](i)
		if err != nil {
			var zero T
			return zero, err
		}
		d2, err := do.Invoke[D2](i)
		if err != nil {
			var zero T
			return zero, err
		}
		return fn(d1, d2), nil
	})
}

// Provide3T registers a typed provider with 3 dependencies.
func Provide3T[T, D1, D2, D3 any](c *Container, fn func(D1, D2, D3) T) {
	do.Provide(c.injector, func(i do.Injector) (T, error) {
		d1, err := do.Invoke[D1](i)
		if err != nil {
			var zero T
			return zero, err
		}
		d2, err := do.Invoke[D2](i)
		if err != nil {
			var zero T
			return zero, err
		}
		d3, err := do.Invoke[D3](i)
		if err != nil {
			var zero T
			return zero, err
		}
		return fn(d1, d2, d3), nil
	})
}

// ProvideValueT provides a pre-created value.
func ProvideValueT[T any](c *Container, value T) {
	do.ProvideValue(c.injector, value)
}

// ResolveAs resolves a dependency with full type safety.
func ResolveAs[T any](c *Container) (T, error) {
	return do.Invoke[T](c.injector)
}

// MustResolveAs resolves a dependency and panics on error.
func MustResolveAs[T any](c *Container) T {
	result, err := ResolveAs[T](c)
	if err != nil {
		panic(fmt.Sprintf("failed to resolve dependency: %v", err))
	}
	return result
}

// Backward compatibility types
type Definition struct {
	Name       string
	Kind       DefinitionKind
	Value      any
	Provider   any
	ModuleName string
	Lazy       bool
	Transient  bool
}

type DefinitionKind string

const (
	DefinitionValue    DefinitionKind = "value"
	DefinitionProvider DefinitionKind = "provider"
)

func (c *Container) Register(def Definition) error {
	switch def.Kind {
	case DefinitionValue:
		if def.Name != "" {
			do.ProvideNamedValue(c.injector, def.Name, def.Value)
		} else {
			do.ProvideValue(c.injector, def.Value)
		}
		return nil
	case DefinitionProvider:
		return nil
	default:
		return fmt.Errorf("unknown definition kind: %v", def.Kind)
	}
}

func (c *Container) Resolve(target any) error {
	return fmt.Errorf("use ResolveAs[T]() for type-safe resolution")
}
