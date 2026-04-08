package dix

import (
	"github.com/samber/do/v2"
	"github.com/samber/oops"
)

// Definition describes a backward-compatible container registration.
type Definition struct {
	Name       string
	Kind       DefinitionKind
	Value      any
	Provider   any
	ModuleName string
	Lazy       bool
	Transient  bool
}

// DefinitionKind describes the kind of backward-compatible registration.
type DefinitionKind string

const (
	// DefinitionValue registers an already-constructed value.
	DefinitionValue DefinitionKind = "value"
	// DefinitionProvider registers a provider function.
	DefinitionProvider DefinitionKind = "provider"
)

// Register registers a backward-compatible definition.
func (c *Container) Register(def Definition) error {
	if c == nil || c.injector == nil {
		return oops.In("dix").
			With("op", "register_legacy_definition", "name", def.Name, "kind", def.Kind, "module", def.ModuleName, "lazy", def.Lazy, "transient", def.Transient).
			New("container is nil")
	}
	switch def.Kind {
	case DefinitionValue:
		if def.Name != "" {
			do.ProvideNamedValue(c.injector, def.Name, def.Value)
		} else {
			do.ProvideValue(c.injector, def.Value)
		}
		return nil
	case DefinitionProvider:
		return oops.In("dix").
			With("op", "register_legacy_definition", "name", def.Name, "kind", def.Kind, "module", def.ModuleName, "lazy", def.Lazy, "transient", def.Transient).
			Errorf("provider definition registration is not implemented; use typed ProviderN helpers instead")
	default:
		return oops.In("dix").
			With("op", "register_legacy_definition", "name", def.Name, "kind", def.Kind, "module", def.ModuleName, "lazy", def.Lazy, "transient", def.Transient).
			Errorf("unknown definition kind: %v", def.Kind)
	}
}

// Resolve keeps backward compatibility for legacy resolve(target) calls.
func (c *Container) Resolve(any) error {
	return oops.In("dix").
		With("op", "resolve_legacy").
		Errorf("resolve(target) is not supported; use ResolveAs[T]() for type-safe resolution")
}
