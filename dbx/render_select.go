package dbx

import (
	"errors"
	"fmt"
	"strings"

	"github.com/DaiYuANg/arcgo/dbx/dialect"
)

func (q *SelectQuery) Build(d dialect.Dialect) (BoundQuery, error) {
	if q == nil {
		return BoundQuery{}, errors.New("dbx: select query is nil")
	}
	if q.FromItem.Name() == "" {
		return BoundQuery{}, errors.New("dbx: select query requires FROM")
	}
	if len(q.Items) == 0 {
		return BoundQuery{}, errors.New("dbx: select query requires at least one item")
	}

	state := &renderState{dialect: d, args: make([]any, 0, 8)}
	if err := renderSelectStatement(state, q); err != nil {
		return BoundQuery{}, err
	}
	if err := state.err(); err != nil {
		return BoundQuery{}, err
	}
	bound := state.BoundQuery()
	if q.LimitN != nil && *q.LimitN > 0 {
		bound.CapacityHint = *q.LimitN
	}
	return bound, nil
}

func renderSelectStatement(state *renderState, q *SelectQuery) error {
	if err := renderCTEs(state, q.CTEs); err != nil {
		return err
	}
	return renderSelectSet(state, q)
}

func renderSelectSet(state *renderState, q *SelectQuery) error {
	if len(q.Unions) == 0 {
		return renderSelectQuery(state, q)
	}

	if err := renderSelectQueryWithoutTail(state, q); err != nil {
		return err
	}
	for _, union := range q.Unions {
		if union.Query == nil {
			return errors.New("dbx: union query is nil")
		}
		if union.All {
			state.writeString(" UNION ALL ")
		} else {
			state.writeString(" UNION ")
		}
		if err := renderUnionQuery(state, union.Query); err != nil {
			return err
		}
	}
	return renderSelectTail(state, q)
}

func renderCTEs(state *renderState, ctes []CTE) error {
	if len(ctes) == 0 {
		return nil
	}
	state.writeString("WITH ")
	for index, cte := range ctes {
		if strings.TrimSpace(cte.Name) == "" {
			return errors.New("dbx: cte name cannot be empty")
		}
		if cte.Query == nil {
			return fmt.Errorf("dbx: cte %s requires query", cte.Name)
		}
		if index > 0 {
			state.writeString(", ")
		}
		state.writeQuotedIdent(strings.TrimSpace(cte.Name))
		state.writeString(" AS (")
		if err := renderSelectStatement(state, cte.Query); err != nil {
			return err
		}
		state.writeByte(')')
	}
	state.writeByte(' ')
	return nil
}

func renderUnionQuery(state *renderState, q *SelectQuery) error {
	if len(q.CTEs) > 0 || len(q.Unions) > 0 || len(q.Orders) > 0 || q.LimitN != nil || q.OffsetN != nil {
		state.writeByte('(')
		if err := renderSelectStatement(state, q); err != nil {
			return err
		}
		state.writeByte(')')
		return nil
	}
	return renderSelectQueryWithoutTail(state, q)
}

func renderSelectQuery(state *renderState, q *SelectQuery) error {
	if err := renderSelectQueryWithoutTail(state, q); err != nil {
		return err
	}
	return renderSelectTail(state, q)
}

func renderSelectQueryWithoutTail(state *renderState, q *SelectQuery) error {
	if err := renderSelectDistinct(state, q); err != nil {
		return err
	}
	if err := renderSelectItems(state, q); err != nil {
		return err
	}
	if err := renderSelectFrom(state, q); err != nil {
		return err
	}
	if err := renderSelectJoins(state, q); err != nil {
		return err
	}
	if err := renderSelectWhere(state, q); err != nil {
		return err
	}
	if err := renderSelectGroupBy(state, q); err != nil {
		return err
	}
	return renderSelectHaving(state, q)
}

func renderSelectDistinct(state *renderState, q *SelectQuery) error {
	state.writeString("SELECT ")
	if q.Distinct {
		state.writeString("DISTINCT ")
	}
	return nil
}

func renderSelectItems(state *renderState, q *SelectQuery) error {
	for i, item := range q.Items {
		if i > 0 {
			state.writeString(", ")
		}
		if err := renderSelectItem(state, item); err != nil {
			return err
		}
	}
	return nil
}

func renderSelectFrom(state *renderState, q *SelectQuery) error {
	state.writeString(" FROM ")
	state.renderTable(q.FromItem)
	return nil
}

func renderSelectJoins(state *renderState, q *SelectQuery) error {
	for _, join := range q.Joins {
		state.writeByte(' ')
		state.writeString(string(join.Type))
		state.writeString(" JOIN ")
		state.renderTable(join.Table)
		if join.Predicate == nil {
			continue
		}
		state.writeString(" ON ")
		if err := renderPredicate(state, join.Predicate); err != nil {
			return err
		}
	}
	return nil
}

func renderSelectWhere(state *renderState, q *SelectQuery) error {
	if q.WhereExp == nil {
		return nil
	}
	state.writeString(" WHERE ")
	return renderPredicate(state, q.WhereExp)
}

func renderSelectGroupBy(state *renderState, q *SelectQuery) error {
	if len(q.Groups) == 0 {
		return nil
	}
	state.writeString(" GROUP BY ")
	for i, group := range q.Groups {
		if i > 0 {
			state.writeString(", ")
		}
		operand, err := renderOperandValue(state, group)
		if err != nil {
			return err
		}
		state.writeString(operand)
	}
	return nil
}

func renderSelectHaving(state *renderState, q *SelectQuery) error {
	if q.HavingExp == nil {
		return nil
	}
	state.writeString(" HAVING ")
	return renderPredicate(state, q.HavingExp)
}

func renderSelectTail(state *renderState, q *SelectQuery) error {
	if err := renderSelectOrders(state, q.Orders); err != nil {
		return err
	}
	return renderSelectLimitOffset(state, q)
}

func renderSelectOrders(state *renderState, orders []Order) error {
	if len(orders) == 0 {
		return nil
	}
	state.writeString(" ORDER BY ")
	for i, order := range orders {
		if i > 0 {
			state.writeString(", ")
		}
		if err := renderOrder(state, order); err != nil {
			return err
		}
	}
	return nil
}

func renderSelectLimitOffset(state *renderState, q *SelectQuery) error {
	clause, err := state.dialect.RenderLimitOffset(q.LimitN, q.OffsetN)
	if err != nil {
		return fmt.Errorf("dbx: render limit offset: %w", err)
	}
	if clause == "" {
		return nil
	}
	state.writeByte(' ')
	state.writeString(clause)
	return nil
}
