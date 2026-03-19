package dialect

type Dialect interface {
	Name() string
	BindVar(n int) string
	QuoteIdent(ident string) string
	RenderLimitOffset(limit, offset *int) (string, error)
}
