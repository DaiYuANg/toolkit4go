package render

// Result contains rendered SQL text and its bind arguments.
type Result struct {
	Query string
	Args  []any
}
