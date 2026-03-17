package dix

import (
	do "github.com/samber/do/v2"
)

// ProvideNamedT registers a named lazy service without dependencies.
func ProvideNamedT[T any](c *Container, name string, fn func() T) {
	do.ProvideNamed(c.injector, name, func(i do.Injector) (T, error) { return fn(), nil })
}

// ProvideNamed1T registers a named lazy service with 1 dependency.
func ProvideNamed1T[T, D1 any](c *Container, name string, fn func(D1) T) {
	do.ProvideNamed(c.injector, name, func(i do.Injector) (T, error) {
		d1, err := do.Invoke[D1](i)
		if err != nil {
			var zero T
			return zero, err
		}
		return fn(d1), nil
	})
}

// ProvideNamed2T registers a named lazy service with 2 dependencies.
func ProvideNamed2T[T, D1, D2 any](c *Container, name string, fn func(D1, D2) T) {
	do.ProvideNamed(c.injector, name, func(i do.Injector) (T, error) {
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

// ProvideNamedValueT registers a named value.
func ProvideNamedValueT[T any](c *Container, name string, value T) {
	do.ProvideNamedValue(c.injector, name, value)
}

// ResolveNamedAs resolves a service by its explicit name.
func ResolveNamedAs[T any](c *Container, name string) (T, error) {
	return do.InvokeNamed[T](c.injector, name)
}

// MustResolveNamedAs resolves a named service or panics.
func MustResolveNamedAs[T any](c *Container, name string) T {
	return do.MustInvokeNamed[T](c.injector, name)
}

// ResolveAssignableAs resolves the first registered service assignable to the requested interface.
// This wraps do.InvokeAs and is the preferred way to use interface-oriented consumption.
func ResolveAssignableAs[T any](c *Container) (T, error) {
	return do.InvokeAs[T](c.injector)
}

// MustResolveAssignableAs resolves an assignable service or panics.
func MustResolveAssignableAs[T any](c *Container) T {
	return do.MustInvokeAs[T](c.injector)
}

// BindAlias binds a concrete service to an interface alias explicitly.
// Use sparingly. For most production cases, ResolveAssignableAs is the better default.
func BindAlias[From, To any](c *Container) error {
	return do.As[From, To](c.injector)
}

// MustBindAlias binds a concrete service to an interface alias or panics.
func MustBindAlias[From, To any](c *Container) {
	do.MustAs[From, To](c.injector)
}

// BindNamedAlias binds a concrete named service to a named interface alias explicitly.
// initial is the existing named concrete service; alias is the new named interface alias.
func BindNamedAlias[From, To any](c *Container, initial string, alias string) error {
	return do.AsNamed[From, To](c.injector, initial, alias)
}

// MustBindNamedAlias binds a concrete named service to a named interface alias or panics.
func MustBindNamedAlias[From, To any](c *Container, initial string, alias string) {
	do.MustAsNamed[From, To](c.injector, initial, alias)
}
