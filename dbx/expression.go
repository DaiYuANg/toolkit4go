package dbx

import "github.com/samber/lo"

type Expression interface {
	expressionNode()
}

type Predicate interface {
	Expression
	predicateNode()
}

type Condition = Predicate

type SelectItem interface {
	selectItemNode()
}

type Assignment interface {
	assignmentNode()
}

type Order interface {
	orderNode()
}

type ComparisonOperator string

type LogicalOperator string

type JoinType string

const (
	OpEq    ComparisonOperator = "="
	OpNe    ComparisonOperator = "<>"
	OpGt    ComparisonOperator = ">"
	OpGe    ComparisonOperator = ">="
	OpLt    ComparisonOperator = "<"
	OpLe    ComparisonOperator = "<="
	OpIn    ComparisonOperator = "IN"
	OpLike  ComparisonOperator = "LIKE"
	OpIs    ComparisonOperator = "IS"
	OpIsNot ComparisonOperator = "IS NOT"
)

const (
	LogicalAnd LogicalOperator = "AND"
	LogicalOr  LogicalOperator = "OR"
)

const (
	InnerJoin JoinType = "INNER"
	LeftJoin  JoinType = "LEFT"
	RightJoin JoinType = "RIGHT"
)

type valueOperand[T any] struct {
	Value T
}

type columnOperand[T any] struct {
	Column typedColumn[T]
}

type comparisonPredicate[E any, T any] struct {
	Left  Column[E, T]
	Op    ComparisonOperator
	Right any
}

func (comparisonPredicate[E, T]) expressionNode() {}
func (comparisonPredicate[E, T]) predicateNode()  {}

type logicalPredicate struct {
	Op         LogicalOperator
	Predicates []Predicate
}

func (logicalPredicate) expressionNode() {}
func (logicalPredicate) predicateNode()  {}

type notPredicate struct {
	Predicate Predicate
}

func (notPredicate) expressionNode() {}
func (notPredicate) predicateNode()  {}

type columnAssignment[E any, T any] struct {
	Column Column[E, T]
	Value  any
}

func (columnAssignment[E, T]) assignmentNode() {}

type columnOrder[E any, T any] struct {
	Column     Column[E, T]
	Descending bool
}

func (columnOrder[E, T]) orderNode() {}

func And(predicates ...Predicate) Predicate {
	items := compactPredicates(predicates)
	if len(items) == 1 {
		return items[0]
	}
	return logicalPredicate{Op: LogicalAnd, Predicates: items}
}

func Or(predicates ...Predicate) Predicate {
	items := compactPredicates(predicates)
	if len(items) == 1 {
		return items[0]
	}
	return logicalPredicate{Op: LogicalOr, Predicates: items}
}

func Not(predicate Predicate) Predicate {
	return notPredicate{Predicate: predicate}
}

func Like[E any](column Column[E, string], pattern string) Predicate {
	return comparisonPredicate[E, string]{
		Left:  column,
		Op:    OpLike,
		Right: valueOperand[string]{Value: pattern},
	}
}

func compactPredicates(predicates []Predicate) []Predicate {
	return lo.Filter(predicates, func(predicate Predicate, _ int) bool {
		return predicate != nil
	})
}
