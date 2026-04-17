package dbx

import (
	"github.com/DaiYuANg/arcgo/collectionx"
	"github.com/DaiYuANg/arcgo/dbx/querydsl"
)

type Expression interface {
	expressionNode()
}

type scalarExpression interface {
	Expression
	operandRenderer
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

type caseWhenBranch struct {
	Predicate Predicate
	Value     any
}

type valueOperand[T any] struct {
	Value T
}

type columnOperand[T any] struct {
	Column typedColumn[T]
}

type excludedColumnOperand[T any] struct {
	Column ColumnMeta
}

type comparisonPredicate struct {
	Left  scalarExpression
	Op    querydsl.ComparisonOperator
	Right any
}

func (comparisonPredicate) expressionNode() {}
func (comparisonPredicate) predicateNode()  {}

type logicalPredicate struct {
	Op         querydsl.LogicalOperator
	Predicates collectionx.List[Predicate]
}

func (logicalPredicate) expressionNode() {}
func (logicalPredicate) predicateNode()  {}

type notPredicate struct {
	Predicate Predicate
}

func (notPredicate) expressionNode() {}
func (notPredicate) predicateNode()  {}

type existsPredicate struct {
	Query *SelectQuery
}

func (existsPredicate) expressionNode() {}
func (existsPredicate) predicateNode()  {}

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

type expressionOrder struct {
	Expr       scalarExpression
	Descending bool
}

func (expressionOrder) orderNode() {}

func And(predicates ...Predicate) Predicate {
	return AndList(compactPredicates(predicates))
}

func Or(predicates ...Predicate) Predicate {
	return OrList(compactPredicates(predicates))
}

func AndList(predicates collectionx.List[Predicate]) Predicate {
	items := compactPredicatesList(predicates)
	if items.Len() == 1 {
		predicate, _ := items.GetFirst()
		return predicate
	}
	return logicalPredicate{Op: querydsl.LogicalAnd, Predicates: items}
}

func OrList(predicates collectionx.List[Predicate]) Predicate {
	items := compactPredicatesList(predicates)
	if items.Len() == 1 {
		predicate, _ := items.GetFirst()
		return predicate
	}
	return logicalPredicate{Op: querydsl.LogicalOr, Predicates: items}
}

func Not(predicate Predicate) Predicate {
	return notPredicate{Predicate: predicate}
}

func Like[E any](column Column[E, string], pattern string) Predicate {
	return comparisonPredicate{
		Left:  column,
		Op:    querydsl.OpLike,
		Right: valueOperand[string]{Value: pattern},
	}
}

func Exists(query *SelectQuery) Predicate {
	return existsPredicate{Query: query}
}

func compactPredicates(predicates []Predicate) collectionx.List[Predicate] {
	return compactPredicatesList(collectionx.NewList(predicates...))
}

func compactPredicatesList(predicates collectionx.List[Predicate]) collectionx.List[Predicate] {
	return collectionx.FilterList(predicates, func(_ int, predicate Predicate) bool {
		return predicate != nil
	})
}

func compactExpressions(expressions []Expression) collectionx.List[Expression] {
	return collectionx.FilterList(collectionx.NewList(expressions...), func(_ int, expression Expression) bool {
		return expression != nil
	})
}
