package repository

import "github.com/DaiYuANg/arcgo/dbx"

// Option configures repository construction behavior.
type Option func(*baseOptions)

type baseOptions struct {
	byIDNotFoundAsError bool
}

func defaultOptions() baseOptions { return baseOptions{} }

// WithByIDNotFoundAsError makes ID-based updates and deletes return ErrNotFound
// when no rows are affected.
func WithByIDNotFoundAsError(enabled bool) Option {
	return func(opts *baseOptions) { opts.byIDNotFoundAsError = enabled }
}

// New constructs a repository with default options.
func New[E any, S EntitySchema[E]](db *dbx.DB, schema S) *Base[E, S] {
	return NewWithOptions[E](db, schema)
}

// NewWithOptions constructs a repository with explicit options.
func NewWithOptions[E any, S EntitySchema[E]](db *dbx.DB, schema S, opts ...Option) *Base[E, S] {
	config := defaultOptions()
	for _, opt := range opts {
		if opt != nil {
			opt(&config)
		}
	}
	return &Base[E, S]{
		db:                  db,
		session:             db,
		schema:              schema,
		mapper:              dbx.MustMapper[E](schema),
		byIDNotFoundAsError: config.byIDNotFoundAsError,
	}
}
