package sqltmplx

import (
	"fmt"

	"github.com/DaiYuANg/arcgo/dbx/dialect"
	"github.com/samber/lo"
)

// Engine compiles and renders SQL templates for a dialect.
type Engine struct {
	dialect dialect.Contract
	cfg     config
}

// New returns a template engine configured for the provided dialect.
func New(d dialect.Contract, opts ...Option) *Engine {
	cfg := config{}
	lo.ForEach(opts, func(opt Option, _ int) {
		if opt != nil {
			opt(&cfg)
		}
	})
	return &Engine{dialect: d, cfg: cfg}
}

// Compile compiles an unnamed template.
func (e *Engine) Compile(tpl string) (*Template, error) {
	return e.CompileNamed("", tpl)
}

// CompileNamed compiles a named template.
func (e *Engine) CompileNamed(name, tpl string) (*Template, error) {
	return compileTemplate(name, tpl, e.dialect, e.cfg)
}

// Render compiles and renders a template with the provided parameters.
func (e *Engine) Render(tpl string, params any) (BoundSQL, error) {
	t, err := e.Compile(tpl)
	if err != nil {
		return BoundSQL{}, fmt.Errorf("compile template: %w", err)
	}
	return t.Render(params)
}
