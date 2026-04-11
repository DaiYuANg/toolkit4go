package sqltmplx

import (
	"github.com/DaiYuANg/arcgo/dbx"
	"github.com/DaiYuANg/arcgo/dbx/sqltmplx/render"
)

// PageRequest is the shared offset-pagination request model.
type PageRequest = dbx.PageRequest

// PageResult contains the items and metadata for a paginated query.
type PageResult[E any] = dbx.PageResult[E]

// Pagination is the normalized template-facing pagination payload.
type Pagination struct {
	Page     int `db:"page"      json:"page"      sqltmpl:"page"`
	PageSize int `db:"page_size" json:"page_size" sqltmpl:"page_size"`
	Limit    int `db:"limit"     json:"limit"     sqltmpl:"limit"`
	Offset   int `db:"offset"    json:"offset"    sqltmpl:"offset"`
}

// Page creates a normalized page request.
func Page(page, pageSize int) PageRequest {
	return dbx.Page(page, pageSize)
}

// NewPageRequest creates a normalized page request.
func NewPageRequest(page, pageSize int) PageRequest {
	return dbx.NewPageRequest(page, pageSize)
}

// NewPagination creates the template-facing pagination payload.
func NewPagination(request PageRequest) Pagination {
	request = request.Normalize()
	return Pagination{
		Page:     request.Page,
		PageSize: request.PageSize,
		Limit:    request.PageSize,
		Offset:   request.Offset(),
	}
}

// WithPage overlays normalized pagination fields under the Page template parameter.
func WithPage(params any, request PageRequest) any {
	return render.WithParam(params, "Page", NewPagination(request))
}

// RenderPage renders the template with normalized pagination parameters.
func (t *Template) RenderPage(params any, request PageRequest) (BoundSQL, error) {
	return t.Render(WithPage(params, request))
}

// BindPage renders the template with normalized pagination parameters into a dbx bound query.
func (t *Template) BindPage(params any, request PageRequest) (dbx.BoundQuery, error) {
	request = request.Normalize()
	bound, err := t.Bind(WithPage(params, request))
	if err != nil {
		return dbx.BoundQuery{}, err
	}
	bound.CapacityHint = request.PageSize
	return bound, nil
}

// RenderPage compiles and renders a template with normalized pagination parameters.
func (e *Engine) RenderPage(tpl string, params any, request PageRequest) (BoundSQL, error) {
	return e.Render(tpl, WithPage(params, request))
}
