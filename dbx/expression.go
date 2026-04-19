package dbx

import (
	"github.com/DaiYuANg/arcgo/collectionx"
	"github.com/DaiYuANg/arcgo/dbx/querydsl"
	schemax "github.com/DaiYuANg/arcgo/dbx/schema"
)

type columnOperand[T any] struct {
	Column typedColumn[T]
}

type excludedColumnOperand[T any] struct {
	Column schemax.ColumnMeta
}

type columnAssignment[E any, T any] struct {
	Column Column[E, T]
	Value  any
}

type columnOrder[E any, T any] struct {
	Column     Column[E, T]
	Descending bool
}

func And(predicates ...querydsl.Predicate) querydsl.Predicate {
	return querydsl.And(predicates...)
}

func Or(predicates ...querydsl.Predicate) querydsl.Predicate {
	return querydsl.Or(predicates...)
}

func AndList(predicates collectionx.List[querydsl.Predicate]) querydsl.Predicate {
	return querydsl.AndList(predicates)
}

func OrList(predicates collectionx.List[querydsl.Predicate]) querydsl.Predicate {
	return querydsl.OrList(predicates)
}

func Not(predicate querydsl.Predicate) querydsl.Predicate {
	return querydsl.Not(predicate)
}

func Like[E any](column Column[E, string], pattern string) querydsl.Predicate {
	return querydsl.Like(column, pattern)
}

func Exists(query *querydsl.SelectQuery) querydsl.Predicate {
	return querydsl.Exists(query)
}

func Select(items ...querydsl.SelectItem) *querydsl.SelectQuery {
	return querydsl.Select(items...)
}

func SelectList(items collectionx.List[querydsl.SelectItem]) *querydsl.SelectQuery {
	return querydsl.SelectList(items)
}

func InsertInto(source querydsl.TableSource) *querydsl.InsertQuery {
	return querydsl.InsertInto(source)
}

func Update(source querydsl.TableSource) *querydsl.UpdateQuery {
	return querydsl.Update(source)
}

func DeleteFrom(source querydsl.TableSource) *querydsl.DeleteQuery {
	return querydsl.DeleteFrom(source)
}

func mergeList[T any](current, next collectionx.List[T]) collectionx.List[T] {
	if current == nil {
		return next.Clone()
	}
	current.Merge(next)
	return current
}
