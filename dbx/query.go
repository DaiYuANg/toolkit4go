package dbx

import "github.com/samber/lo"

type Join struct {
	Type      JoinType
	Table     Table
	Predicate Predicate
}

type SelectQuery struct {
	Items    []SelectItem
	FromItem Table
	Joins    []Join
	WhereExp Predicate
	Orders   []Order
	LimitN   *int
	OffsetN  *int
	Distinct bool
}

type InsertQuery struct {
	Into        Table
	Assignments []Assignment
}

type UpdateQuery struct {
	Table       Table
	Assignments []Assignment
	WhereExp    Predicate
}

type DeleteQuery struct {
	From     Table
	WhereExp Predicate
}

type JoinBuilder struct {
	query *SelectQuery
	index int
}

func Select(items ...SelectItem) *SelectQuery {
	return &SelectQuery{Items: compactSelectItems(items)}
}

func (q *SelectQuery) WithDistinct() *SelectQuery {
	q.Distinct = true
	return q
}

func (q *SelectQuery) From(source TableSource) *SelectQuery {
	q.FromItem = source.tableRef()
	return q
}

func (q *SelectQuery) Where(predicate Predicate) *SelectQuery {
	q.WhereExp = predicate
	return q
}

func (q *SelectQuery) OrderBy(orders ...Order) *SelectQuery {
	q.Orders = append(q.Orders, compactOrders(orders)...)
	return q
}

func (q *SelectQuery) Limit(limit int) *SelectQuery {
	q.LimitN = &limit
	return q
}

func (q *SelectQuery) Offset(offset int) *SelectQuery {
	q.OffsetN = &offset
	return q
}

func (q *SelectQuery) Join(source TableSource) *JoinBuilder {
	q.Joins = append(q.Joins, Join{Type: InnerJoin, Table: source.tableRef()})
	return &JoinBuilder{query: q, index: len(q.Joins) - 1}
}

func (q *SelectQuery) LeftJoin(source TableSource) *JoinBuilder {
	q.Joins = append(q.Joins, Join{Type: LeftJoin, Table: source.tableRef()})
	return &JoinBuilder{query: q, index: len(q.Joins) - 1}
}

func (q *SelectQuery) RightJoin(source TableSource) *JoinBuilder {
	q.Joins = append(q.Joins, Join{Type: RightJoin, Table: source.tableRef()})
	return &JoinBuilder{query: q, index: len(q.Joins) - 1}
}

func (b *JoinBuilder) On(predicate Predicate) *SelectQuery {
	b.query.Joins[b.index].Predicate = predicate
	return b.query
}

func InsertInto(source TableSource) *InsertQuery {
	return &InsertQuery{Into: source.tableRef()}
}

func (q *InsertQuery) Values(assignments ...Assignment) *InsertQuery {
	q.Assignments = append(q.Assignments, compactAssignments(assignments)...)
	return q
}

func Update(source TableSource) *UpdateQuery {
	return &UpdateQuery{Table: source.tableRef()}
}

func (q *UpdateQuery) Set(assignments ...Assignment) *UpdateQuery {
	q.Assignments = append(q.Assignments, compactAssignments(assignments)...)
	return q
}

func (q *UpdateQuery) Where(predicate Predicate) *UpdateQuery {
	q.WhereExp = predicate
	return q
}

func DeleteFrom(source TableSource) *DeleteQuery {
	return &DeleteQuery{From: source.tableRef()}
}

func (q *DeleteQuery) Where(predicate Predicate) *DeleteQuery {
	q.WhereExp = predicate
	return q
}

func compactSelectItems(items []SelectItem) []SelectItem {
	return lo.Filter(items, func(item SelectItem, _ int) bool {
		return item != nil
	})
}

func compactOrders(orders []Order) []Order {
	return lo.Filter(orders, func(order Order, _ int) bool {
		return order != nil
	})
}

func compactAssignments(assignments []Assignment) []Assignment {
	return lo.Filter(assignments, func(assignment Assignment, _ int) bool {
		return assignment != nil
	})
}
