package advanced

import (
	collectionmap "github.com/DaiYuANg/arcgo/collectionx/mapping"
	"github.com/DaiYuANg/arcgo/dix"
	"github.com/samber/do/v2"
	"github.com/samber/lo"
)

// Inspection summarizes advanced runtime inspection output.
type Inspection struct {
	ScopeTree         string
	ProvidedServices  []do.ServiceDescription
	InvokedServices   []do.ServiceDescription
	NamedDependencies map[string]string
}

// InspectOptions controls which inspection sections are populated.
type InspectOptions struct {
	IncludeScopeTree        bool
	IncludeProvidedServices bool
	IncludeInvokedServices  bool
	IncludeNamedDeps        bool
}

// DefaultInspectOptions returns the default inspection option set.
func DefaultInspectOptions() InspectOptions {
	return InspectOptions{
		IncludeScopeTree:        true,
		IncludeProvidedServices: true,
		IncludeInvokedServices:  true,
		IncludeNamedDeps:        true,
	}
}

// ExplainScopeTree returns the textual do scope tree for a runtime.
func ExplainScopeTree(rt *dix.Runtime) string {
	if rt == nil {
		return ""
	}

	explainedScope := do.ExplainInjector(rt.Raw())
	return explainedScope.String()
}

// ListProvidedServices returns the services provided by the runtime injector.
func ListProvidedServices(rt *dix.Runtime) []do.ServiceDescription {
	if rt == nil {
		return nil
	}

	return rt.Raw().ListProvidedServices()
}

// ListInvokedServices returns the services invoked by the runtime injector.
func ListInvokedServices(rt *dix.Runtime) []do.ServiceDescription {
	if rt == nil {
		return nil
	}

	return rt.Raw().ListInvokedServices()
}

// ExplainNamedDependencies returns dependency trees for the requested named services.
func ExplainNamedDependencies(rt *dix.Runtime, namedServices ...string) map[string]string {
	if rt == nil || len(namedServices) == 0 {
		return nil
	}

	dependencies := collectionmap.NewMapWithCapacity[string, string](len(namedServices))
	lo.ForEach(namedServices, func(name string, _ int) {
		if desc, found := do.ExplainNamedService(rt.Raw(), name); found {
			dependencies.Set(name, desc.String())
		}
	})

	return dependencies.All()
}

// InspectRuntime inspects a runtime with the default options.
func InspectRuntime(rt *dix.Runtime, namedServices ...string) Inspection {
	return InspectRuntimeWithOptions(rt, DefaultInspectOptions(), namedServices...)
}

// InspectRuntimeWithOptions inspects a runtime with the provided options.
func InspectRuntimeWithOptions(rt *dix.Runtime, opts InspectOptions, namedServices ...string) Inspection {
	if rt == nil {
		return Inspection{}
	}

	var scopeTree string
	if opts.IncludeScopeTree {
		scopeTree = ExplainScopeTree(rt)
	}

	var provided []do.ServiceDescription
	if opts.IncludeProvidedServices {
		provided = ListProvidedServices(rt)
	}

	var invoked []do.ServiceDescription
	if opts.IncludeInvokedServices {
		invoked = ListInvokedServices(rt)
	}

	var namedDependencies map[string]string
	if opts.IncludeNamedDeps && len(namedServices) > 0 {
		namedDependencies = ExplainNamedDependencies(rt, namedServices...)
	}

	return Inspection{
		ScopeTree:         scopeTree,
		ProvidedServices:  provided,
		InvokedServices:   invoked,
		NamedDependencies: namedDependencies,
	}
}
