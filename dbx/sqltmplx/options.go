package sqltmplx

import "github.com/DaiYuANg/arcgo/dbx/sqltmplx/validate"

// Option configures the template engine.
type Option func(*config)

type config struct {
	validator validate.Validator
}

// WithValidator configures SQL validation for rendered templates.
func WithValidator(v validate.Validator) Option {
	return func(c *config) {
		c.validator = v
	}
}
