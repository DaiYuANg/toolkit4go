package repository

import "github.com/DaiYuANg/arcgo/dbx"

type Base[E any, S dbx.SchemaSource[E]] struct {
	db     *dbx.DB
	schema S
	mapper dbx.Mapper[E]
}

func New[E any, S dbx.SchemaSource[E]](db *dbx.DB, schema S) *Base[E, S] {
	return &Base[E, S]{db: db, schema: schema, mapper: dbx.MustMapper[E](schema)}
}

func (r *Base[E, S]) DB() *dbx.DB {
	return r.db
}

func (r *Base[E, S]) Schema() S {
	return r.schema
}

func (r *Base[E, S]) Mapper() dbx.Mapper[E] {
	return r.mapper
}
