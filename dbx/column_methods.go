package dbx

import (
	"reflect"

	"github.com/samber/lo"
)

func (c Column[E, T]) expressionNode() {}
func (c Column[E, T]) selectItemNode() {}
func (c Column[E, T]) refNode()        {}

func (Column[E, T]) valueType() reflect.Type {
	return reflect.TypeFor[T]()
}

func (c Column[E, T]) Name() string {
	return c.meta.Name
}

func (c Column[E, T]) TableName() string {
	return c.meta.Table
}

func (c Column[E, T]) TableAlias() string {
	return c.meta.Alias
}

func (c Column[E, T]) FieldName() string {
	return c.meta.FieldName
}

func (c Column[E, T]) Meta() ColumnMeta {
	meta := c.meta
	if meta.References != nil {
		meta.References = new(*meta.References)
	}
	return meta
}

func (c Column[E, T]) IsPrimaryKey() bool {
	return c.meta.PrimaryKey
}

func (c Column[E, T]) IsNullable() bool {
	return c.meta.Nullable
}

func (c Column[E, T]) IsUnique() bool {
	return c.meta.Unique
}

func (c Column[E, T]) IsIndexed() bool {
	return c.meta.Indexed
}

func (c Column[E, T]) DefaultValue() string {
	return c.meta.DefaultValue
}

func (c Column[E, T]) Reference() (ForeignKeyRef, bool) {
	if c.meta.References == nil {
		return ForeignKeyRef{}, false
	}
	return *c.meta.References, true
}

func (c Column[E, T]) Ref() string {
	if c.meta.Alias != "" {
		return c.meta.Alias + "." + c.meta.Name
	}
	return c.meta.Table + "." + c.meta.Name
}

func (c Column[E, T]) Eq(value T) Predicate {
	return comparisonPredicate{
		Left:  c,
		Op:    OpEq,
		Right: valueOperand[T]{Value: value},
	}
}

func (c Column[E, T]) EqColumn(other typedColumn[T]) Predicate {
	return comparisonPredicate{
		Left:  c,
		Op:    OpEq,
		Right: columnOperand[T]{Column: other},
	}
}

func (c Column[E, T]) Ne(value T) Predicate {
	return comparisonPredicate{
		Left:  c,
		Op:    OpNe,
		Right: valueOperand[T]{Value: value},
	}
}

func (c Column[E, T]) Gt(value T) Predicate {
	return comparisonPredicate{
		Left:  c,
		Op:    OpGt,
		Right: valueOperand[T]{Value: value},
	}
}

func (c Column[E, T]) Ge(value T) Predicate {
	return comparisonPredicate{
		Left:  c,
		Op:    OpGe,
		Right: valueOperand[T]{Value: value},
	}
}

func (c Column[E, T]) Lt(value T) Predicate {
	return comparisonPredicate{
		Left:  c,
		Op:    OpLt,
		Right: valueOperand[T]{Value: value},
	}
}

func (c Column[E, T]) Le(value T) Predicate {
	return comparisonPredicate{
		Left:  c,
		Op:    OpLe,
		Right: valueOperand[T]{Value: value},
	}
}

func (c Column[E, T]) In(values ...T) Predicate {
	return comparisonPredicate{
		Left: c,
		Op:   OpIn,
		Right: lo.Map(values, func(value T, _ int) any {
			return value
		}),
	}
}

func (c Column[E, T]) InQuery(query *SelectQuery) Predicate {
	return comparisonPredicate{
		Left:  c,
		Op:    OpIn,
		Right: subqueryOperand{Query: query},
	}
}

func (c Column[E, T]) IsNull() Predicate {
	return comparisonPredicate{
		Left: c,
		Op:   OpIs,
	}
}

func (c Column[E, T]) IsNotNull() Predicate {
	return comparisonPredicate{
		Left: c,
		Op:   OpIsNot,
	}
}

func (c Column[E, T]) Set(value T) Assignment {
	return columnAssignment[E, T]{
		Column: c,
		Value:  valueOperand[T]{Value: value},
	}
}

func (c Column[E, T]) SetColumn(other typedColumn[T]) Assignment {
	return columnAssignment[E, T]{
		Column: c,
		Value:  columnOperand[T]{Column: other},
	}
}

func (c Column[E, T]) SetExcluded() Assignment {
	return columnAssignment[E, T]{
		Column: c,
		Value:  excludedColumnOperand[T]{Column: c.columnRef()},
	}
}

func (c Column[E, T]) Asc() Order {
	return columnOrder[E, T]{Column: c}
}

func (c Column[E, T]) Desc() Order {
	return columnOrder[E, T]{Column: c, Descending: true}
}

func (c Column[E, T]) columnRef() ColumnMeta {
	return c.meta
}

func (c Column[E, T]) As(alias string) SelectItem {
	return aliasedSelectItem{Item: c, Alias: alias}
}
