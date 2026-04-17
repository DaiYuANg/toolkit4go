package dix

import (
	"github.com/samber/mo"
	"github.com/samber/oops"
)

// ResolveAs resolves a typed value from the container.
func ResolveAs[T any](c *Container) (T, error) {
	if c == nil || c.injector == nil {
		var zero T
		return zero, oops.In("dix").
			With("op", "resolve", "service", serviceNameOf[T]()).
			New("container is nil")
	}
	return resolveInjectorAs[T](c.injector)
}

// ResolveOptionalAs resolves an optional typed value from the container.
func ResolveOptionalAs[T any](c *Container) (value T, ok bool) {
	return ResolveOptional[T](c)
}

// ResolveOptional resolves an optional typed value from the container.
func ResolveOptional[T any](c *Container) (value T, ok bool) {
	return ResolveOption[T](c).Get()
}

// ResolveOptionAs resolves an optional dependency as mo.Option.
func ResolveOptionAs[T any](c *Container) mo.Option[T] {
	return ResolveOption[T](c)
}

// ResolveOption resolves an optional dependency as mo.Option.
func ResolveOption[T any](c *Container) mo.Option[T] {
	value, err := ResolveAs[T](c)
	if err == nil {
		return mo.Some(value)
	}
	return mo.None[T]()
}

// ResolveOrElse resolves a typed value or returns the provided fallback.
func ResolveOrElse[T any](c *Container, fallback T) T {
	return ResolveOr[T](c, fallback)
}

// ResolveOr resolves a typed value or returns the provided fallback.
func ResolveOr[T any](c *Container, fallback T) T {
	return ResolveOption[T](c).OrElse(fallback)
}

// MustResolveAs resolves a typed value and panics on failure.
func MustResolveAs[T any](c *Container) T {
	result, err := ResolveAs[T](c)
	if err != nil {
		panic(oops.In("dix").
			With("op", "must_resolve", "service", serviceNameOf[T]()).
			Wrapf(err, "resolve dependency"))
	}
	return result
}
