package dix

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	collectionlist "github.com/DaiYuANg/arcgo/collectionx/list"
)

type buildPlan struct {
	spec    *appSpec
	modules *collectionlist.List[*moduleSpec]
}

func newBuildPlan(app *App) (*buildPlan, error) {
	if app == nil || app.spec == nil {
		return nil, fmt.Errorf("app is nil")
	}

	modules, err := flattenModuleList(&app.spec.modules, app.spec.profile)
	if err != nil {
		logger := app.spec.logger
		if logger != nil {
			logger.Error("module flatten failed", "app", app.Name(), "error", err)
		}
		return nil, fmt.Errorf("module flatten failed: %w", err)
	}

	plan := &buildPlan{
		spec:    app.spec,
		modules: modules,
	}

	return plan, validateTypedGraph(plan)
}

func (p *buildPlan) Build() (*Runtime, error) {
	if p == nil || p.spec == nil {
		return nil, fmt.Errorf("build plan is nil")
	}

	logger := p.spec.logger
	if logger == nil {
		logger = slog.Default()
	}
	debugEnabled := logger.Enabled(context.Background(), slog.LevelDebug)
	infoEnabled := logger.Enabled(context.Background(), slog.LevelInfo)

	if infoEnabled {
		logger.Info("building app", "app", p.spec.meta.Name, "profile", p.spec.profile)
	}
	if debugEnabled {
		logger.Debug("build plan ready",
			"app", p.spec.meta.Name,
			"modules", p.modules.Len(),
			"providers", countModuleProviders(p.modules),
			"hooks", countModuleHooks(p.modules),
			"setups", countModuleSetups(p.modules),
			"invokes", countModuleInvokes(p.modules),
		)
	}

	rt := newRuntime(p.spec, p)
	ProvideValueT[*slog.Logger](rt.container, rt.logger)
	ProvideValueT[AppMeta](rt.container, p.spec.meta)
	ProvideValueT[Profile](rt.container, p.spec.profile)

	p.modules.Range(func(_ int, mod *moduleSpec) bool {
		if debugEnabled {
			logger.Debug("registering module",
				"module", mod.name,
				"providers", mod.providers.Len(),
				"hooks", mod.hooks.Len(),
				"setups", mod.setups.Len(),
				"invokes", mod.invokes.Len(),
			)
		}
		mod.providers.Range(func(_ int, provider ProviderFunc) bool {
			if debugEnabled {
				logger.Debug("registering provider",
					"module", mod.name,
					"label", provider.meta.Label,
					"output", provider.meta.Output.Name,
					"dependencies", serviceRefNames(provider.meta.Dependencies),
					"raw", provider.meta.Raw,
				)
			}
			provider.apply(rt.container)
			return true
		})
		return true
	})

	var setupErr error
	p.modules.Range(func(_ int, mod *moduleSpec) bool {
		mod.hooks.Range(func(_ int, hook HookFunc) bool {
			if debugEnabled {
				logger.Debug("binding lifecycle hook",
					"module", mod.name,
					"label", hook.meta.Label,
					"kind", hook.meta.Kind,
					"dependencies", serviceRefNames(hook.meta.Dependencies),
					"raw", hook.meta.Raw,
				)
			}
			hook.bind(rt.container, rt.lifecycle)
			return true
		})

		mod.setups.Range(func(_ int, setup SetupFunc) bool {
			if debugEnabled {
				logger.Debug("running module setup",
					"module", mod.name,
					"label", setup.meta.Label,
					"dependencies", serviceRefNames(setup.meta.Dependencies),
					"provides", serviceRefNames(setup.meta.Provides),
					"overrides", serviceRefNames(setup.meta.Overrides),
					"graph_mutation", setup.meta.GraphMutation,
					"raw", setup.meta.Raw,
				)
			}
			if err := setup.apply(rt.container, rt.lifecycle); err != nil {
				logger.Error("module setup failed", "module", mod.name, "label", setup.meta.Label, "error", err)
				setupErr = fmt.Errorf("setup failed for module %s via %s: %w", mod.name, setup.meta.Label, err)
				return false
			}
			if debugEnabled {
				logger.Debug("module setup completed", "module", mod.name, "label", setup.meta.Label)
			}
			return true
		})
		return setupErr == nil
	})
	if setupErr != nil {
		return nil, cleanupBuildFailure(rt, logger, setupErr)
	}

	var buildErr error
	p.modules.Range(func(_ int, mod *moduleSpec) bool {
		var invokeErr error
		mod.invokes.Range(func(_ int, invoke InvokeFunc) bool {
			if debugEnabled {
				logger.Debug("running invoke",
					"module", mod.name,
					"label", invoke.meta.Label,
					"dependencies", serviceRefNames(invoke.meta.Dependencies),
					"raw", invoke.meta.Raw,
				)
			}
			invokeErr = invoke.apply(rt.container)
			if invokeErr == nil && debugEnabled {
				logger.Debug("invoke completed", "module", mod.name, "label", invoke.meta.Label)
			}
			return invokeErr == nil
		})
		if invokeErr != nil {
			logger.Error("invoke failed", "module", mod.name, "error", invokeErr)
			buildErr = fmt.Errorf("invoke failed in module %s: %w", mod.name, invokeErr)
			return false
		}
		return true
	})
	if buildErr != nil {
		return nil, cleanupBuildFailure(rt, logger, buildErr)
	}

	rt.state = AppStateBuilt
	if infoEnabled {
		logger.Info("app built", "app", rt.Name(), "modules", p.modules.Len())
	}
	rt.logDebugInformation()
	return rt, nil
}

func cleanupBuildFailure(rt *Runtime, logger *slog.Logger, buildErr error) error {
	if rt == nil || rt.container == nil {
		return buildErr
	}

	report := rt.container.ShutdownReport(context.Background())
	if report == nil || len(report.Errors) == 0 {
		return buildErr
	}
	if logger != nil {
		logger.Error("build cleanup failed", "app", rt.Name(), "error", report)
	}
	return errors.Join(buildErr, report)
}

func countModuleProviders(modules *collectionlist.List[*moduleSpec]) int {
	total := 0
	if modules == nil {
		return total
	}
	modules.Range(func(_ int, mod *moduleSpec) bool {
		if mod != nil {
			total += mod.providers.Len()
		}
		return true
	})
	return total
}

func countModuleHooks(modules *collectionlist.List[*moduleSpec]) int {
	total := 0
	if modules == nil {
		return total
	}
	modules.Range(func(_ int, mod *moduleSpec) bool {
		if mod != nil {
			total += mod.hooks.Len()
		}
		return true
	})
	return total
}

func countModuleSetups(modules *collectionlist.List[*moduleSpec]) int {
	total := 0
	if modules == nil {
		return total
	}
	modules.Range(func(_ int, mod *moduleSpec) bool {
		if mod != nil {
			total += mod.setups.Len()
		}
		return true
	})
	return total
}

func countModuleInvokes(modules *collectionlist.List[*moduleSpec]) int {
	total := 0
	if modules == nil {
		return total
	}
	modules.Range(func(_ int, mod *moduleSpec) bool {
		if mod != nil {
			total += mod.invokes.Len()
		}
		return true
	})
	return total
}

func serviceRefNames(refs []ServiceRef) []string {
	names := make([]string, 0, len(refs))
	for _, ref := range refs {
		if ref.Name != "" {
			names = append(names, ref.Name)
		}
	}
	return names
}
