package advanced

import (
	"github.com/DaiYuANg/arcgo/dix"
	"github.com/samber/do/v2"
)

// ResolveInjectorAs resolves a typed value directly from a do injector.
func ResolveInjectorAs[T any](injector do.Injector) (T, error) {
	return do.InvokeNamed[T](injector, typedName[T]())
}

// MustResolveInjectorAs resolves a typed value directly from a do injector and panics on failure.
func MustResolveInjectorAs[T any](injector do.Injector) T {
	return do.MustInvokeNamed[T](injector, typedName[T]())
}

// ResolveRuntimeAs resolves a typed value from a dix runtime.
func ResolveRuntimeAs[T any](rt *dix.Runtime) (T, error) {
	if rt == nil {
		var zero T
		return zero, do.ErrServiceNotFound
	}
	return ResolveInjectorAs[T](rt.Raw())
}

// MustResolveRuntimeAs resolves a typed value from a dix runtime and panics on failure.
func MustResolveRuntimeAs[T any](rt *dix.Runtime) T {
	return MustResolveInjectorAs[T](rt.Raw())
}

// ResolveNamedAs resolves a named value from a dix container.
func ResolveNamedAs[T any](c *dix.Container, name string) (T, error) {
	return do.InvokeNamed[T](c.Raw(), name)
}

// MustResolveNamedAs resolves a named value from a dix container and panics on failure.
func MustResolveNamedAs[T any](c *dix.Container, name string) T {
	return do.MustInvokeNamed[T](c.Raw(), name)
}

// ResolveAssignableAs resolves an assignable value from a dix container.
func ResolveAssignableAs[T any](c *dix.Container) (T, error) {
	return do.InvokeAs[T](c.Raw())
}

// MustResolveAssignableAs resolves an assignable value from a dix container and panics on failure.
func MustResolveAssignableAs[T any](c *dix.Container) T {
	return do.MustInvokeAs[T](c.Raw())
}
