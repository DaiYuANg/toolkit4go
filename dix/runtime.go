package dix

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	do "github.com/samber/do/v2"
)

func newRuntime(spec *appSpec, plan *buildPlan) *Runtime {
	logger := defaultLogger()
	if spec != nil && spec.logger != nil {
		logger = spec.logger
	}

	rt := &Runtime{
		spec:      spec,
		plan:      plan,
		container: newContainer(),
		lifecycle: newLifecycle(),
		logger:    logger,
		state:     AppStateCreated,
	}

	rt.container.logger = rt.logger
	rt.lifecycle.logger = rt.logger

	return rt
}

func (r *Runtime) Name() string {
	if r == nil || r.spec == nil {
		return ""
	}
	return r.spec.meta.Name
}

func (r *Runtime) Profile() Profile {
	if r == nil || r.spec == nil {
		return ""
	}
	return r.spec.profile
}

func (r *Runtime) Logger() *slog.Logger {
	if r == nil {
		return nil
	}
	return r.logger
}

func (r *Runtime) Meta() AppMeta {
	if r == nil || r.spec == nil {
		return AppMeta{}
	}
	return r.spec.meta
}

func (r *Runtime) State() AppState {
	if r == nil {
		return AppStateCreated
	}
	return r.state
}

func (r *Runtime) Container() *Container {
	if r == nil {
		return nil
	}
	return r.container
}

func (r *Runtime) Raw() do.Injector {
	if r == nil || r.container == nil {
		return nil
	}
	return r.container.Raw()
}

func (r *Runtime) Start(ctx context.Context) error {
	if r == nil {
		return fmt.Errorf("runtime is nil")
	}
	if r.state != AppStateBuilt {
		return fmt.Errorf("runtime must be built before starting")
	}

	r.transitionState(AppStateStarting, "start requested")
	if r.logger.Enabled(context.Background(), slog.LevelInfo) {
		r.logger.Info("starting app", "app", r.Name())
	}
	startedHooks, err := r.lifecycle.executeStartHooks(ctx, r.container)
	if err != nil {
		if r.logger.Enabled(context.Background(), slog.LevelDebug) {
			r.logger.Debug("rolling back app start",
				"app", r.Name(),
				"started_hooks", startedHooks,
				"rollback_stop_hooks", startedHooks,
			)
		}
		rollbackErr := r.lifecycle.executeStopHooksSubset(ctx, startedHooks)
		shutdownReport := r.container.ShutdownReport(ctx)
		startErr := errors.Join(err, rollbackErr, shutdownReport)
		r.transitionState(AppStateStopped, "start failed")
		r.logger.Error("app start failed", "app", r.Name(), "error", startErr)
		return startErr
	}

	r.transitionState(AppStateStarted, "start completed")
	if r.logger.Enabled(context.Background(), slog.LevelInfo) {
		r.logger.Info("app started", "app", r.Name())
	}
	return nil
}

func (r *Runtime) Stop(ctx context.Context) error {
	report, err := r.StopWithReport(ctx)
	if err != nil {
		return err
	}
	if report != nil {
		return report.Err()
	}
	return nil
}

func (r *Runtime) StopWithReport(ctx context.Context) (*StopReport, error) {
	if r == nil {
		return nil, fmt.Errorf("runtime is nil")
	}
	if r.state == AppStateStarting {
		return nil, fmt.Errorf("runtime is still starting")
	}
	if r.state != AppStateStarted {
		return nil, fmt.Errorf("runtime must be started before stopping")
	}

	if r.logger.Enabled(context.Background(), slog.LevelInfo) {
		r.logger.Info("stopping app", "app", r.Name())
	}
	if r.logger.Enabled(context.Background(), slog.LevelDebug) {
		r.logger.Debug("executing runtime stop",
			"app", r.Name(),
			"stop_hooks", len(r.lifecycle.stopHooks.Values()),
		)
	}

	report := &StopReport{}
	if err := r.lifecycle.executeStopHooks(ctx, r.container); err != nil {
		r.logger.Error("stop hooks failed", "app", r.Name(), "error", err)
		report.HookError = err
	}

	report.ShutdownReport = r.container.ShutdownReport(ctx)
	if report.ShutdownReport != nil && len(report.ShutdownReport.Errors) > 0 {
		r.logger.Error("container shutdown failed", "app", r.Name(), "error", report.ShutdownReport)
	}
	if r.logger.Enabled(context.Background(), slog.LevelDebug) {
		shutdownErrors := 0
		if report.ShutdownReport != nil {
			shutdownErrors = len(report.ShutdownReport.Errors)
		}
		r.logger.Debug("runtime stop report",
			"app", r.Name(),
			"hook_error", report.HookError != nil,
			"shutdown_errors", shutdownErrors,
		)
	}

	r.transitionState(AppStateStopped, "stop completed")
	if r.logger.Enabled(context.Background(), slog.LevelInfo) {
		r.logger.Info("app stopped", "app", r.Name())
	}

	return report, report.Err()
}

func (r *Runtime) logDebugInformation() {
	if r == nil || r.spec == nil {
		return
	}

	if r.spec.debug.scopeTree {
		injector := do.ExplainInjector(r.container.Raw())
		r.logger.Info("do scope tree", "app", r.Name(), "tree", injector.String())
	}

	r.spec.debug.namedServiceDependencies.Range(func(name string) bool {
		if desc, found := do.ExplainNamedService(r.container.Raw(), name); found {
			r.logger.Info("do named service dependencies", "app", r.Name(), "name", name, "dependencies", desc.String())
		} else {
			r.logger.Warn("do named service not found", "app", r.Name(), "name", name)
		}
		return true
	})
}

func (r *Runtime) transitionState(next AppState, reason string) {
	if r == nil {
		return
	}
	prev := r.state
	r.state = next
	if r.logger != nil && r.logger.Enabled(context.Background(), slog.LevelDebug) && prev != next {
		r.logger.Debug("runtime state transition",
			"app", r.Name(),
			"from", prev.String(),
			"to", next.String(),
			"reason", reason,
		)
	}
}
