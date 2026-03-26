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

	known := collectionset.NewSetWithCapacity[string](64,
		serviceNameOf[*slog.Logger](),
		serviceNameOf[AppMeta](),
		serviceNameOf[Profile](),
	)
	errs := collectionlist.NewListWithCapacity[error](4)
	warnings := collectionlist.NewListWithCapacity[ValidationWarning](2)

	validateDeps := func(modName, kind, label string, deps []ServiceRef) {
		validateDependencies(errs, known, modName, kind, label, deps)
	}

	addWarning := func(kind ValidationWarningKind, moduleName, label, details string) {
		warnings.Add(ValidationWarning{
			Kind:    kind,
			Module:  moduleName,
			Label:   label,
			Details: details,
		})
	}

	plan.modules.Range(func(_ int, mod *moduleSpec) bool {
		if mod == nil {
			return true
		}
		modName := mod.name

		mod.providers.Range(func(_ int, provider ProviderFunc) bool {
			meta := provider.meta
			if meta.Output.Name != "" {
				if known.Contains(meta.Output.Name) {
					errs.Add(fmt.Errorf("duplicate provider output `%s` in module `%s` via %s", meta.Output.Name, modName, meta.Label))
					return true
				}
				known.Add(meta.Output.Name)
			} else if meta.Raw {
				addWarning(
					ValidationWarningRawProviderUndeclaredOutput,
					modName,
					meta.Label,
					"raw provider has no declared output; validation cannot model services it registers",
				)
			}
			return true
		})

		mod.setups.Range(func(_ int, setup SetupFunc) bool {
			meta := setup.meta
			for _, provide := range meta.Provides {
				if known.Contains(provide.Name) {
					errs.Add(fmt.Errorf("duplicate setup output `%s` in module `%s` via %s", provide.Name, modName, meta.Label))
					continue
				}
				known.Add(provide.Name)
			}
			if meta.Raw && len(meta.Provides) == 0 && len(meta.Overrides) == 0 && meta.GraphMutation {
				addWarning(
					ValidationWarningRawSetupUndeclaredGraph,
					modName,
					meta.Label,
					"raw setup mutates the graph without declared provides/overrides; validation cannot model its graph effects",
				)
			}
			return true
		})

		return true
	})

	plan.modules.Range(func(_ int, mod *moduleSpec) bool {
		if mod == nil {
			return true
		}
		modName := mod.name

		mod.providers.Range(func(_ int, provider ProviderFunc) bool {
			meta := provider.meta
			if meta.Raw && len(meta.Dependencies) == 0 {
				addWarning(
					ValidationWarningRawProviderUndeclaredDeps,
					modName,
					meta.Label,
					"raw provider has no declared dependencies; validation cannot verify what it resolves at registration time",
				)
			}
			if len(meta.Dependencies) == 0 {
				return true
			}
			validateDeps(modName, "provider", meta.Label, meta.Dependencies)
			return true
		})

		mod.setups.Range(func(_ int, setup SetupFunc) bool {
			meta := setup.meta
			for _, override := range meta.Overrides {
				if !known.Contains(override.Name) {
					errs.Add(fmt.Errorf("override target `%s` not found in module `%s` via %s", override.Name, modName, meta.Label))
					continue
				}
			}
			validateDeps(modName, "setup", meta.Label, meta.Dependencies)
			return true
		})

		mod.invokes.Range(func(_ int, invoke InvokeFunc) bool {
			if invoke.meta.Raw && len(invoke.meta.Dependencies) == 0 {
				addWarning(
					ValidationWarningRawInvokeUndeclaredDeps,
					modName,
					invoke.meta.Label,
					"raw invoke has no declared dependencies; validation cannot verify what it resolves",
				)
				return true
			}
			if len(invoke.meta.Dependencies) == 0 {
				return true
			}
			validateDeps(modName, "invoke", invoke.meta.Label, invoke.meta.Dependencies)
			return true
		})

		mod.hooks.Range(func(_ int, hook HookFunc) bool {
			if hook.meta.Raw && len(hook.meta.Dependencies) == 0 {
				addWarning(
					ValidationWarningRawHookUndeclaredDeps,
					modName,
					hook.meta.Label,
					"raw hook has no declared dependencies; validation cannot verify what it resolves during lifecycle execution",
				)
				return true
			}
			if len(hook.meta.Dependencies) == 0 {
				return true
			}
			validateDeps(modName, string(hook.meta.Kind)+" hook", hook.meta.Label, hook.meta.Dependencies)
			return true
		})

		return true
	})

	return ValidationReport{
		Errors:   errs.Values(),
		Warnings: warnings.Values(),
	}
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
