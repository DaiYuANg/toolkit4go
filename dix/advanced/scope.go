package advanced

import (
	"github.com/DaiYuANg/arcgo/dix"
	"github.com/samber/do/v2"
	"github.com/samber/lo"
)

// ScopePackage configures a newly created do scope.
type ScopePackage func(do.Injector)

// Scope creates a named child scope from a runtime injector.
func Scope(rt *dix.Runtime, name string, packages ...ScopePackage) *do.Scope {
	if rt == nil {
		return nil
	}

	valid := lo.Filter(packages, func(pkg ScopePackage, _ int) bool { return pkg != nil })
	switch len(valid) {
	case 0:
		return rt.Raw().Scope(name)
	case 1:
		current := valid[0]
		return rt.Raw().Scope(name, func(injector do.Injector) {
			current(injector)
		})
	default:
		wrapped := lo.Map(valid, func(pkg ScopePackage, _ int) func(do.Injector) {
			current := pkg
			return func(injector do.Injector) {
				current(injector)
			}
		})
		return rt.Raw().Scope(name, wrapped...)
	}
}

// ProvideScopedValue registers a typed singleton value in a scope.
func ProvideScopedValue[T any](injector do.Injector, value T) {
	do.ProvideNamedValue(injector, typedName[T](), value)
}

// ProvideScopedNamedValue registers a named singleton value in a scope.
func ProvideScopedNamedValue[T any](injector do.Injector, name string, value T) {
	do.ProvideNamedValue(injector, name, value)
}

// ProvideScoped0 registers a typed scoped provider with no dependencies.
func ProvideScoped0[T any](injector do.Injector, fn func() T) {
	do.ProvideNamed(injector, typedName[T](), func(do.Injector) (T, error) {
		return fn(), nil
	})
}

// ProvideScopedNamed0 registers a named scoped provider with no dependencies.
func ProvideScopedNamed0[T any](injector do.Injector, name string, fn func() T) {
	do.ProvideNamed(injector, name, func(do.Injector) (T, error) {
		return fn(), nil
	})
}

// ProvideScoped1 registers a typed scoped provider with one dependency.
func ProvideScoped1[T, D1 any](injector do.Injector, fn func(D1) T) {
	do.ProvideNamed(injector, typedName[T](), func(i do.Injector) (T, error) {
		d1, err := invokeTyped[D1](i)
		if err != nil {
			var zero T
			return zero, err
		}
		return fn(d1), nil
	})
}

// ProvideScopedNamed1 registers a named scoped provider with one dependency.
func ProvideScopedNamed1[T, D1 any](injector do.Injector, name string, fn func(D1) T) {
	do.ProvideNamed(injector, name, func(i do.Injector) (T, error) {
		d1, err := invokeTyped[D1](i)
		if err != nil {
			var zero T
			return zero, err
		}
		return fn(d1), nil
	})
}

// ProvideScoped2 registers a typed scoped provider with two dependencies.
func ProvideScoped2[T, D1, D2 any](injector do.Injector, fn func(D1, D2) T) {
	do.ProvideNamed(injector, typedName[T](), func(i do.Injector) (T, error) {
		d1, err := invokeTyped[D1](i)
		if err != nil {
			var zero T
			return zero, err
		}
		d2, err := invokeTyped[D2](i)
		if err != nil {
			var zero T
			return zero, err
		}
		return fn(d1, d2), nil
	})
}

// ProvideScopedNamed2 registers a named scoped provider with two dependencies.
func ProvideScopedNamed2[T, D1, D2 any](injector do.Injector, name string, fn func(D1, D2) T) {
	do.ProvideNamed(injector, name, func(i do.Injector) (T, error) {
		d1, err := invokeTyped[D1](i)
		if err != nil {
			var zero T
			return zero, err
		}
		d2, err := invokeTyped[D2](i)
		if err != nil {
			var zero T
			return zero, err
		}
		return fn(d1, d2), nil
	})
}

// ProvideScoped3 registers a typed scoped provider with three dependencies.
func ProvideScoped3[T, D1, D2, D3 any](injector do.Injector, fn func(D1, D2, D3) T) {
	do.ProvideNamed(injector, typedName[T](), func(i do.Injector) (T, error) {
		d1, err := invokeTyped[D1](i)
		if err != nil {
			var zero T
			return zero, err
		}
		d2, err := invokeTyped[D2](i)
		if err != nil {
			var zero T
			return zero, err
		}
		d3, err := invokeTyped[D3](i)
		if err != nil {
			var zero T
			return zero, err
		}
		return fn(d1, d2, d3), nil
	})
}

// ProvideScopedNamed3 registers a named scoped provider with three dependencies.
func ProvideScopedNamed3[T, D1, D2, D3 any](injector do.Injector, name string, fn func(D1, D2, D3) T) {
	do.ProvideNamed(injector, name, func(i do.Injector) (T, error) {
		d1, err := invokeTyped[D1](i)
		if err != nil {
			var zero T
			return zero, err
		}
		d2, err := invokeTyped[D2](i)
		if err != nil {
			var zero T
			return zero, err
		}
		d3, err := invokeTyped[D3](i)
		if err != nil {
			var zero T
			return zero, err
		}
		return fn(d1, d2, d3), nil
	})
}

// ResolveScopedAs resolves a typed value from a scope injector.
func ResolveScopedAs[T any](injector do.Injector) (T, error) {
	return ResolveInjectorAs[T](injector)
}

// ResolveScopedNamedAs resolves a named value from a scope injector.
func ResolveScopedNamedAs[T any](injector do.Injector, name string) (T, error) {
	return do.InvokeNamed[T](injector, name)
}
