package dix

import (
	"errors"
	"strings"
)

// Validate validates the immutable app spec and current module graph.
func (a *App) Validate() error {
	return a.ValidateReport().Err()
}

// ValidateReport validates the app and returns the full validation report.
func (a *App) ValidateReport() ValidationReport {
	plan, err := newBuildPlan(a)
	if err != nil {
		return ValidationReport{Errors: []error{err}}
	}
	return validateTypedGraphReport(plan)
}

// HasWarnings reports whether the validation report contains warnings.
func (r ValidationReport) HasWarnings() bool {
	return len(r.Warnings) > 0
}

// HasErrors reports whether the validation report contains errors.
func (r ValidationReport) HasErrors() bool {
	return len(r.Errors) > 0
}

// Err returns the combined validation error.
func (r ValidationReport) Err() error {
	return errors.Join(r.Errors...)
}

// WarningSummary renders the validation warnings as a newline-delimited summary.
func (r ValidationReport) WarningSummary() string {
	if len(r.Warnings) == 0 {
		return ""
	}

	lines := make([]string, 0, len(r.Warnings))
	for _, warning := range r.Warnings {
		line := string(warning.Kind)
		if warning.Module != "" {
			line += " module=" + warning.Module
		}
		if warning.Label != "" {
			line += " label=" + warning.Label
		}
		if warning.Details != "" {
			line += " " + warning.Details
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}
