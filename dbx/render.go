package dbx

import (
	"fmt"
	"strings"

	"github.com/DaiYuANg/arcgo/dbx/dialect"
)

type predicateRenderer interface {
	renderPredicate(*renderState) error
}

type assignmentRenderer interface {
	renderAssignment(*renderState) error
}

type insertAssignmentRenderer interface {
	assignmentRenderer
	renderAssignmentValue(*renderState) error
	assignmentColumn() ColumnMeta
}

type orderRenderer interface {
	renderOrder(*renderState) error
}

type operandRenderer interface {
	renderOperand(*renderState) (string, error)
}

type renderState struct {
	dialect dialect.Dialect
	buf     strings.Builder
	args    []any
}

func (s *renderState) bind(value any) string {
	s.args = append(s.args, value)
	return s.dialect.BindVar(len(s.args))
}

func (s *renderState) writeQuotedIdent(name string) {
	s.buf.WriteString(s.dialect.QuoteIdent(name))
}

func (s *renderState) writeQualifiedIdent(table, column string) {
	if table != "" {
		s.writeQuotedIdent(table)
		s.buf.WriteByte('.')
	}
	s.writeQuotedIdent(column)
}

func (s *renderState) renderColumn(meta ColumnMeta) {
	table := meta.Table
	if meta.Alias != "" {
		table = meta.Alias
	}
	s.writeQualifiedIdent(table, meta.Name)
}

func (s *renderState) renderTable(table Table) {
	s.writeQuotedIdent(table.Name())
	if alias := table.Alias(); alias != "" && alias != table.Name() {
		s.buf.WriteString(" AS ")
		s.writeQuotedIdent(alias)
	}
}

func (s *renderState) BoundQuery() BoundQuery {
	args := make([]any, len(s.args))
	copy(args, s.args)
	return BoundQuery{SQL: s.buf.String(), Args: args}
}

func (q *SelectQuery) Build(d dialect.Dialect) (BoundQuery, error) {
	if q == nil {
		return BoundQuery{}, fmt.Errorf("dbx: select query is nil")
	}
	if q.FromItem.Name() == "" {
		return BoundQuery{}, fmt.Errorf("dbx: select query requires FROM")
	}
	if len(q.Items) == 0 {
		return BoundQuery{}, fmt.Errorf("dbx: select query requires at least one item")
	}

	state := &renderState{dialect: d, args: make([]any, 0, 8)}
	state.buf.WriteString("SELECT ")
	if q.Distinct {
		state.buf.WriteString("DISTINCT ")
	}
	for i, item := range q.Items {
		if i > 0 {
			state.buf.WriteString(", ")
		}
		if err := renderSelectItem(state, item); err != nil {
			return BoundQuery{}, err
		}
	}

	state.buf.WriteString(" FROM ")
	state.renderTable(q.FromItem)
	for _, join := range q.Joins {
		state.buf.WriteByte(' ')
		state.buf.WriteString(string(join.Type))
		state.buf.WriteString(" JOIN ")
		state.renderTable(join.Table)
		if join.Predicate != nil {
			state.buf.WriteString(" ON ")
			if err := renderPredicate(state, join.Predicate); err != nil {
				return BoundQuery{}, err
			}
		}
	}
	if q.WhereExp != nil {
		state.buf.WriteString(" WHERE ")
		if err := renderPredicate(state, q.WhereExp); err != nil {
			return BoundQuery{}, err
		}
	}
	if len(q.Orders) > 0 {
		state.buf.WriteString(" ORDER BY ")
		for i, order := range q.Orders {
			if i > 0 {
				state.buf.WriteString(", ")
			}
			if err := renderOrder(state, order); err != nil {
				return BoundQuery{}, err
			}
		}
	}
	clause, err := d.RenderLimitOffset(q.LimitN, q.OffsetN)
	if err != nil {
		return BoundQuery{}, err
	}
	if clause != "" {
		state.buf.WriteByte(' ')
		state.buf.WriteString(clause)
	}
	return state.BoundQuery(), nil
}

func (q *InsertQuery) Build(d dialect.Dialect) (BoundQuery, error) {
	if q == nil {
		return BoundQuery{}, fmt.Errorf("dbx: insert query is nil")
	}
	if q.Into.Name() == "" {
		return BoundQuery{}, fmt.Errorf("dbx: insert query requires target table")
	}
	if len(q.Assignments) == 0 {
		return BoundQuery{}, fmt.Errorf("dbx: insert query requires assignments")
	}

	state := &renderState{dialect: d, args: make([]any, 0, len(q.Assignments))}
	state.buf.WriteString("INSERT INTO ")
	state.renderTable(q.Into)
	state.buf.WriteString(" (")
	for i, assignment := range q.Assignments {
		renderer, ok := assignment.(insertAssignmentRenderer)
		if !ok {
			return BoundQuery{}, fmt.Errorf("dbx: unsupported insert assignment %T", assignment)
		}
		if i > 0 {
			state.buf.WriteString(", ")
		}
		state.writeQuotedIdent(renderer.assignmentColumn().Name)
	}
	state.buf.WriteString(") VALUES (")
	for i, assignment := range q.Assignments {
		renderer, ok := assignment.(insertAssignmentRenderer)
		if !ok {
			return BoundQuery{}, fmt.Errorf("dbx: unsupported insert assignment %T", assignment)
		}
		if i > 0 {
			state.buf.WriteString(", ")
		}
		if err := renderer.renderAssignmentValue(state); err != nil {
			return BoundQuery{}, err
		}
	}
	state.buf.WriteByte(')')
	return state.BoundQuery(), nil
}

func (q *UpdateQuery) Build(d dialect.Dialect) (BoundQuery, error) {
	if q == nil {
		return BoundQuery{}, fmt.Errorf("dbx: update query is nil")
	}
	if q.Table.Name() == "" {
		return BoundQuery{}, fmt.Errorf("dbx: update query requires target table")
	}
	if len(q.Assignments) == 0 {
		return BoundQuery{}, fmt.Errorf("dbx: update query requires assignments")
	}

	state := &renderState{dialect: d, args: make([]any, 0, len(q.Assignments))}
	state.buf.WriteString("UPDATE ")
	state.renderTable(q.Table)
	state.buf.WriteString(" SET ")
	for i, assignment := range q.Assignments {
		if i > 0 {
			state.buf.WriteString(", ")
		}
		if err := renderAssignment(state, assignment); err != nil {
			return BoundQuery{}, err
		}
	}
	if q.WhereExp != nil {
		state.buf.WriteString(" WHERE ")
		if err := renderPredicate(state, q.WhereExp); err != nil {
			return BoundQuery{}, err
		}
	}
	return state.BoundQuery(), nil
}

func (q *DeleteQuery) Build(d dialect.Dialect) (BoundQuery, error) {
	if q == nil {
		return BoundQuery{}, fmt.Errorf("dbx: delete query is nil")
	}
	if q.From.Name() == "" {
		return BoundQuery{}, fmt.Errorf("dbx: delete query requires target table")
	}

	state := &renderState{dialect: d, args: make([]any, 0, 4)}
	state.buf.WriteString("DELETE FROM ")
	state.renderTable(q.From)
	if q.WhereExp != nil {
		state.buf.WriteString(" WHERE ")
		if err := renderPredicate(state, q.WhereExp); err != nil {
			return BoundQuery{}, err
		}
	}
	return state.BoundQuery(), nil
}

func renderSelectItem(state *renderState, item SelectItem) error {
	column, ok := item.(columnAccessor)
	if !ok {
		return fmt.Errorf("dbx: unsupported select item %T", item)
	}
	state.renderColumn(column.columnRef())
	return nil
}

func renderPredicate(state *renderState, predicate Predicate) error {
	renderer, ok := predicate.(predicateRenderer)
	if !ok {
		return fmt.Errorf("dbx: unsupported predicate %T", predicate)
	}
	return renderer.renderPredicate(state)
}

func renderAssignment(state *renderState, assignment Assignment) error {
	renderer, ok := assignment.(assignmentRenderer)
	if !ok {
		return fmt.Errorf("dbx: unsupported assignment %T", assignment)
	}
	return renderer.renderAssignment(state)
}

func renderOrder(state *renderState, order Order) error {
	renderer, ok := order.(orderRenderer)
	if !ok {
		return fmt.Errorf("dbx: unsupported order %T", order)
	}
	return renderer.renderOrder(state)
}

func (o valueOperand[T]) renderOperand(state *renderState) (string, error) {
	return state.bind(o.Value), nil
}

func (o columnOperand[T]) renderOperand(state *renderState) (string, error) {
	meta := o.Column.columnRef()
	var builder strings.Builder
	table := meta.Table
	if meta.Alias != "" {
		table = meta.Alias
	}
	builder.WriteString(state.dialect.QuoteIdent(table))
	builder.WriteByte('.')
	builder.WriteString(state.dialect.QuoteIdent(meta.Name))
	return builder.String(), nil
}

func (p comparisonPredicate[E, T]) renderPredicate(state *renderState) error {
	state.renderColumn(p.Left.columnRef())
	if p.Op == OpIs || p.Op == OpIsNot {
		state.buf.WriteByte(' ')
		state.buf.WriteString(string(p.Op))
		state.buf.WriteString(" NULL")
		return nil
	}
	operand, err := renderOperandValue(state, p.Right)
	if err != nil {
		return err
	}
	state.buf.WriteByte(' ')
	state.buf.WriteString(string(p.Op))
	state.buf.WriteByte(' ')
	state.buf.WriteString(operand)
	return nil
}

func (p logicalPredicate) renderPredicate(state *renderState) error {
	if len(p.Predicates) == 0 {
		return fmt.Errorf("dbx: logical predicate requires nested predicates")
	}
	state.buf.WriteByte('(')
	for i, predicate := range p.Predicates {
		if i > 0 {
			state.buf.WriteByte(' ')
			state.buf.WriteString(string(p.Op))
			state.buf.WriteByte(' ')
		}
		if err := renderPredicate(state, predicate); err != nil {
			return err
		}
	}
	state.buf.WriteByte(')')
	return nil
}

func (p notPredicate) renderPredicate(state *renderState) error {
	if p.Predicate == nil {
		return fmt.Errorf("dbx: NOT predicate requires nested predicate")
	}
	state.buf.WriteString("NOT (")
	if err := renderPredicate(state, p.Predicate); err != nil {
		return err
	}
	state.buf.WriteByte(')')
	return nil
}

func (a columnAssignment[E, T]) assignmentColumn() ColumnMeta {
	return a.Column.columnRef()
}

func (a columnAssignment[E, T]) renderAssignment(state *renderState) error {
	state.writeQuotedIdent(a.Column.Name())
	state.buf.WriteString(" = ")
	operand, err := renderOperandValue(state, a.Value)
	if err != nil {
		return err
	}
	state.buf.WriteString(operand)
	return nil
}

func (a columnAssignment[E, T]) renderAssignmentValue(state *renderState) error {
	operand, err := renderOperandValue(state, a.Value)
	if err != nil {
		return err
	}
	state.buf.WriteString(operand)
	return nil
}

func (o columnOrder[E, T]) renderOrder(state *renderState) error {
	state.renderColumn(o.Column.columnRef())
	if o.Descending {
		state.buf.WriteString(" DESC")
		return nil
	}
	state.buf.WriteString(" ASC")
	return nil
}

func renderOperandValue(state *renderState, value any) (string, error) {
	if renderer, ok := value.(operandRenderer); ok {
		return renderer.renderOperand(state)
	}
	if values, ok := value.([]any); ok {
		return renderAnySliceOperand(state, values)
	}
	return state.bind(value), nil
}

func renderAnySliceOperand(state *renderState, values []any) (string, error) {
	if len(values) == 0 {
		return "", fmt.Errorf("dbx: IN operand cannot be empty")
	}
	var builder strings.Builder
	builder.WriteByte('(')
	for i, value := range values {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(state.bind(value))
	}
	builder.WriteByte(')')
	return builder.String(), nil
}
