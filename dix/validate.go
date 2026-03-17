package dix

import "fmt"

// Validate checks the module graph before runtime execution.
// It detects import cycles and other structural issues early.
func (a *App) Validate() error {
	if _, err := flattenModules(a.modules, a.profile); err != nil {
		return fmt.Errorf("module graph validation failed: %w", err)
	}
	return nil
}
