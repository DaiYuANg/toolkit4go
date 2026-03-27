package sqltmplx

// BoundSQL contains rendered SQL text and its bind arguments.
type BoundSQL struct {
	Query string
	Args  []any
}
