package dbx_test

import (
	"context"
	"fmt"

	dbx "github.com/DaiYuANg/arcgo/dbx"
	"github.com/samber/mo"
)

func LoadBelongsTo[S any, T any](ctx context.Context, session Session, sources []S, sourceSchema SchemaSource[S], sourceMapper Mapper[S], relation BelongsTo[S, T], targetSchema SchemaSource[T], targetMapper Mapper[T], assign func(int, *S, mo.Option[T])) error {
	if err := dbx.LoadBelongsTo(ctx, session, sources, sourceSchema, sourceMapper, relation, targetSchema, targetMapper, assign); err != nil {
		return fmt.Errorf("load belongs-to relation: %w", err)
	}
	return nil
}

func LoadHasMany[S any, T any](ctx context.Context, session Session, sources []S, sourceSchema SchemaSource[S], sourceMapper Mapper[S], relation HasMany[S, T], targetSchema SchemaSource[T], targetMapper Mapper[T], assign func(int, *S, []T)) error {
	if err := dbx.LoadHasMany(ctx, session, sources, sourceSchema, sourceMapper, relation, targetSchema, targetMapper, assign); err != nil {
		return fmt.Errorf("load has-many relation: %w", err)
	}
	return nil
}

func LoadHasOne[S any, T any](ctx context.Context, session Session, sources []S, sourceSchema SchemaSource[S], sourceMapper Mapper[S], relation HasOne[S, T], targetSchema SchemaSource[T], targetMapper Mapper[T], assign func(int, *S, mo.Option[T])) error {
	if err := dbx.LoadHasOne(ctx, session, sources, sourceSchema, sourceMapper, relation, targetSchema, targetMapper, assign); err != nil {
		return fmt.Errorf("load has-one relation: %w", err)
	}
	return nil
}

func LoadManyToMany[S any, T any](ctx context.Context, session Session, sources []S, sourceSchema SchemaSource[S], sourceMapper Mapper[S], relation ManyToMany[S, T], targetSchema SchemaSource[T], targetMapper Mapper[T], assign func(int, *S, []T)) error {
	if err := dbx.LoadManyToMany(ctx, session, sources, sourceSchema, sourceMapper, relation, targetSchema, targetMapper, assign); err != nil {
		return fmt.Errorf("load many-to-many relation: %w", err)
	}
	return nil
}
