package dix

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/samber/lo"
)

// HealthKind is the category of a health check.
type HealthKind string

const (
	HealthKindGeneral   HealthKind = "general"
	HealthKindLiveness  HealthKind = "liveness"
	HealthKindReadiness HealthKind = "readiness"
)

// HealthCheckFunc is a framework-level health check.
type HealthCheckFunc func(context.Context) error

type healthCheckEntry struct {
	name string
	kind HealthKind
	fn   HealthCheckFunc
}

// RegisterHealthCheck registers a framework-managed general health check.
func (c *Container) RegisterHealthCheck(name string, fn HealthCheckFunc) {
	c.RegisterHealthCheckOfKind(HealthKindGeneral, name, fn)
}

// RegisterLivenessCheck registers a liveness health check.
func (c *Container) RegisterLivenessCheck(name string, fn HealthCheckFunc) {
	c.RegisterHealthCheckOfKind(HealthKindLiveness, name, fn)
}

// RegisterReadinessCheck registers a readiness health check.
func (c *Container) RegisterReadinessCheck(name string, fn HealthCheckFunc) {
	c.RegisterHealthCheckOfKind(HealthKindReadiness, name, fn)
}

// RegisterHealthCheckOfKind registers a categorized health check.
func (c *Container) RegisterHealthCheckOfKind(kind HealthKind, name string, fn HealthCheckFunc) {
	if c == nil || fn == nil {
		return
	}
	if kind == "" {
		kind = HealthKindGeneral
	}
	c.healthChecks = append(c.healthChecks, healthCheckEntry{name: name, kind: kind, fn: fn})
}

// HealthReport describes the current health status.
type HealthReport struct {
	Kind   HealthKind       `json:"kind"`
	Checks map[string]error `json:"-"`
}

// Healthy reports whether all checks passed.
func (r HealthReport) Healthy() bool {
	for _, err := range r.Checks {
		if err != nil {
			return false
		}
	}
	return true
}

// Error returns a combined error when one or more checks fail.
func (r HealthReport) Error() error {
	if r.Healthy() {
		return nil
	}

	names := lo.FilterMap(lo.Keys(r.Checks), func(name string, _ int) (string, bool) {
		err := r.Checks[name]
		if err == nil {
			return "", false
		}
		return fmt.Sprintf("%s: %v", name, err), true
	})
	sort.Strings(names)
	return fmt.Errorf("health check failed: %s", strings.Join(names, "; "))
}

// MarshalJSON renders a user-friendly JSON payload for HTTP endpoints.
func (r HealthReport) MarshalJSON() ([]byte, error) {
	type payload struct {
		Kind    HealthKind         `json:"kind"`
		Healthy bool               `json:"healthy"`
		Checks  map[string]*string `json:"checks"`
	}
	checks := make(map[string]*string, len(r.Checks))
	for name, err := range r.Checks {
		if err == nil {
			checks[name] = nil
			continue
		}
		msg := err.Error()
		checks[name] = &msg
	}
	return json.Marshal(payload{Kind: r.Kind, Healthy: r.Healthy(), Checks: checks})
}

// CheckHealth executes all framework-managed health checks.
func (a *App) CheckHealth(ctx context.Context) HealthReport {
	return a.checkHealthByKind(ctx, HealthKindGeneral)
}

// CheckLiveness executes all liveness checks.
func (a *App) CheckLiveness(ctx context.Context) HealthReport {
	return a.checkHealthByKind(ctx, HealthKindLiveness)
}

// CheckReadiness executes all readiness checks.
func (a *App) CheckReadiness(ctx context.Context) HealthReport {
	return a.checkHealthByKind(ctx, HealthKindReadiness)
}

func (a *App) checkHealthByKind(ctx context.Context, kind HealthKind) HealthReport {
	report := HealthReport{Kind: kind, Checks: make(map[string]error)}
	for _, check := range a.container.healthChecks {
		if check.kind != kind {
			continue
		}
		err := check.fn(ctx)
		report.Checks[check.name] = err
		if a.logger != nil {
			if err != nil {
				a.logger.Warn("health check failed", "kind", check.kind, "check", check.name, "error", err)
			} else {
				a.logger.Debug("health check passed", "kind", check.kind, "check", check.name)
			}
		}
	}
	return report
}

func (a *App) healthHandler(kind HealthKind) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var report HealthReport
		switch kind {
		case HealthKindLiveness:
			report = a.CheckLiveness(ctx)
		case HealthKindReadiness:
			report = a.CheckReadiness(ctx)
		default:
			report = a.CheckHealth(ctx)
		}

		status := http.StatusOK
		if !report.Healthy() {
			status = http.StatusServiceUnavailable
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(report)
	}
}
