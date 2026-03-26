package dix

import (
	"fmt"
	"log/slog"

	collectionlist "github.com/DaiYuANg/arcgo/collectionx/list"
	collectionset "github.com/DaiYuANg/arcgo/collectionx/set"
)

func validateTypedGraph(plan *buildPlan) error {
	return validateTypedGraphReport(plan).Err()
}

func validateTypedGraphReport(plan *buildPlan) ValidationReport {
	if plan == nil || plan.modules == nil {
		return ValidationReport{}
	}

	state := newValidationState()
	collectDeclaredOutputs(plan.modules, state)
	validateDeclaredDependencies(plan.modules, state)

	return ValidationReport{
		Errors:   state.errs.Values(),
		Warnings: state.warnings.Values(),
	}
}

type validationState struct {
	known    *collectionset.Set[string]
	errs     *collectionlist.List[error]
	warnings *collectionlist.List[ValidationWarning]
}

func newValidationState() *validationState {
	return &validationState{
		known: collectionset.NewSetWithCapacity[string](64,
			serviceNameOf[*slog.Logger](),
			serviceNameOf[AppMeta](),
			serviceNameOf[Profile](),
		),
		errs:     collectionlist.NewListWithCapacity[error](4),
		warnings: collectionlist.NewListWithCapacity[ValidationWarning](2),
	}
}

func collectDeclaredOutputs(modules *collectionlist.List[*moduleSpec], state *validationState) {
	modules.Range(func(_ int, mod *moduleSpec) bool {
		if mod == nil {
			return true
		}
		collectProviderOutputs(mod, state)
		collectSetupOutputs(mod, state)
		return true
	})
}

func collectProviderOutputs(mod *moduleSpec, state *validationState) {
	mod.providers.Range(func(_ int, provider ProviderFunc) bool {
		meta := provider.meta
		if meta.Output.Name != "" {
			if state.known.Contains(meta.Output.Name) {
				state.errs.Add(fmt.Errorf("duplicate provider output `%s` in module `%s` via %s", meta.Output.Name, mod.name, meta.Label))
				return true
			}
			state.known.Add(meta.Output.Name)
			return true
		}
		if meta.Raw {
			state.addWarning(
				ValidationWarningRawProviderUndeclaredOutput,
				mod.name,
				meta.Label,
				"raw provider has no declared output; validation cannot model services it registers",
			)
		}
		return true
	})
}

func collectSetupOutputs(mod *moduleSpec, state *validationState) {
	mod.setups.Range(func(_ int, setup SetupFunc) bool {
		meta := setup.meta
		for _, provide := range meta.Provides {
			if state.known.Contains(provide.Name) {
				state.errs.Add(fmt.Errorf("duplicate setup output `%s` in module `%s` via %s", provide.Name, mod.name, meta.Label))
				continue
			}
			state.known.Add(provide.Name)
		}
		if meta.Raw && len(meta.Provides) == 0 && len(meta.Overrides) == 0 && meta.GraphMutation {
			state.addWarning(
				ValidationWarningRawSetupUndeclaredGraph,
				mod.name,
				meta.Label,
				"raw setup mutates the graph without declared provides/overrides; validation cannot model its graph effects",
			)
		}
		return true
	})
}

func validateDeclaredDependencies(modules *collectionlist.List[*moduleSpec], state *validationState) {
	modules.Range(func(_ int, mod *moduleSpec) bool {
		if mod == nil {
			return true
		}
		validateProviderDependencies(mod, state)
		validateSetupDependencies(mod, state)
		validateInvokeDependencies(mod, state)
		validateHookDependencies(mod, state)
		return true
	})
}

func validateProviderDependencies(mod *moduleSpec, state *validationState) {
	mod.providers.Range(func(_ int, provider ProviderFunc) bool {
		meta := provider.meta
		if meta.Raw && len(meta.Dependencies) == 0 {
			state.addWarning(
				ValidationWarningRawProviderUndeclaredDeps,
				mod.name,
				meta.Label,
				"raw provider has no declared dependencies; validation cannot verify what it resolves at registration time",
			)
		}
		state.validateDeps(mod.name, "provider", meta.Label, meta.Dependencies)
		return true
	})
}

func validateSetupDependencies(mod *moduleSpec, state *validationState) {
	mod.setups.Range(func(_ int, setup SetupFunc) bool {
		meta := setup.meta
		for _, override := range meta.Overrides {
			if !state.known.Contains(override.Name) {
				state.errs.Add(fmt.Errorf("override target `%s` not found in module `%s` via %s", override.Name, mod.name, meta.Label))
			}
		}
		state.validateDeps(mod.name, "setup", meta.Label, meta.Dependencies)
		return true
	})
}

func validateInvokeDependencies(mod *moduleSpec, state *validationState) {
	mod.invokes.Range(func(_ int, invoke InvokeFunc) bool {
		if invoke.meta.Raw && len(invoke.meta.Dependencies) == 0 {
			state.addWarning(
				ValidationWarningRawInvokeUndeclaredDeps,
				mod.name,
				invoke.meta.Label,
				"raw invoke has no declared dependencies; validation cannot verify what it resolves",
			)
			return true
		}
		state.validateDeps(mod.name, "invoke", invoke.meta.Label, invoke.meta.Dependencies)
		return true
	})
}

func validateHookDependencies(mod *moduleSpec, state *validationState) {
	mod.hooks.Range(func(_ int, hook HookFunc) bool {
		if hook.meta.Raw && len(hook.meta.Dependencies) == 0 {
			state.addWarning(
				ValidationWarningRawHookUndeclaredDeps,
				mod.name,
				hook.meta.Label,
				"raw hook has no declared dependencies; validation cannot verify what it resolves during lifecycle execution",
			)
			return true
		}
		state.validateDeps(mod.name, string(hook.meta.Kind)+" hook", hook.meta.Label, hook.meta.Dependencies)
		return true
	})
}

func (s *validationState) addWarning(kind ValidationWarningKind, moduleName, label, details string) {
	s.warnings.Add(ValidationWarning{
		Kind:    kind,
		Module:  moduleName,
		Label:   label,
		Details: details,
	})
}

func (s *validationState) validateDeps(moduleName, kind, label string, deps []ServiceRef) {
	validateDependencies(s.errs, s.known, moduleName, kind, label, deps)
}

func validateDependencies(
	errs *collectionlist.List[error],
	known *collectionset.Set[string],
	moduleName string,
	kind string,
	label string,
	deps []ServiceRef,
) {
	for _, dep := range deps {
		if !known.Contains(dep.Name) {
			errs.Add(fmt.Errorf("missing dependency `%s` for %s %s in module `%s`", dep.Name, kind, label, moduleName))
		}
	}
}
