package sqltmplx

import (
	"github.com/DaiYuANg/arcgo/dbx"
	"github.com/DaiYuANg/arcgo/dbx/paging"
	"github.com/DaiYuANg/arcgo/dbx/sqltmplx/render"
)

// Page creates a normalized page request.
func Page(page, pageSize int) paging.Request {
	return paging.Page(page, pageSize)
}

// NewPageRequest creates a normalized page request.
func NewPageRequest(page, pageSize int) paging.Request {
	return paging.NewRequest(page, pageSize)
}

// WithPage overlays a normalized paging.Request under the Page template parameter.
func WithPage(params any, request paging.Request) any {
	return render.WithParam(params, "Page", request.Normalize())
}

// RenderPage renders the template with normalized pagination parameters.
func (t *Template) RenderPage(params any, request paging.Request) (BoundSQL, error) {
	return t.Render(WithPage(params, request))
}

// BindPage renders the template with normalized pagination parameters into a dbx bound query.
func (t *Template) BindPage(params any, request paging.Request) (dbx.BoundQuery, error) {
	request = request.Normalize()
	bound, err := t.Bind(WithPage(params, request))
	if err != nil {
		return dbx.BoundQuery{}, err
	}
	bound.CapacityHint = request.Limit()
	return bound, nil
}

// RenderPage compiles and renders a template with normalized pagination parameters.
func (e *Engine) RenderPage(tpl string, params any, request paging.Request) (BoundSQL, error) {
	return e.Render(tpl, WithPage(params, request))
}
