package activerecord

import (
	"github.com/DaiYuANg/arcgo/dbx"
	"github.com/DaiYuANg/arcgo/dbx/repository"
)

type Store[E any, S dbx.SchemaSource[E]] struct {
	repository *repository.Base[E, S]
}

func New[E any, S dbx.SchemaSource[E]](db *dbx.DB, schema S) *Store[E, S] {
	return &Store[E, S]{repository: repository.New[E](db, schema)}
}

func (s *Store[E, S]) Repository() *repository.Base[E, S] {
	return s.repository
}
