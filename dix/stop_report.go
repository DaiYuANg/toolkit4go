package dix

import (
	"errors"
	"strings"

	"github.com/samber/do/v2"
	"github.com/samber/lo"
)

// StopReport aggregates errors produced while stopping a runtime.
type StopReport struct {
	HookError      error
	ShutdownReport *do.ShutdownReport
}

// HasErrors reports whether the stop report contains any errors.
func (r *StopReport) HasErrors() bool {
	return r != nil && r.Err() != nil
}

func (r *StopReport) collectErrors() []error {
	if r == nil {
		return nil
	}
	errs := lo.Compact([]error{r.HookError})
	if r.ShutdownReport != nil && len(r.ShutdownReport.Errors) > 0 {
		errs = append(errs, r.ShutdownReport)
	}
	return errs
}

// Err returns the combined stop error.
func (r *StopReport) Err() error {
	return errors.Join(r.collectErrors()...)
}

// Error returns the combined stop error string.
func (r *StopReport) Error() string {
	errs := r.collectErrors()
	if len(errs) == 0 {
		return ""
	}
	return strings.Join(lo.Map(errs, func(e error, _ int) string { return e.Error() }), "\n")
}
