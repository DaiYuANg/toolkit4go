package dbx

import (
	"errors"
	"fmt"

	"github.com/DaiYuANg/arcgo/collectionx"
	"github.com/DaiYuANg/arcgo/dbx/dialect"
)

func (q *InsertQuery) Build(d dialect.Dialect) (BoundQuery, error) {
	if q == nil {
		return BoundQuery{}, errors.New("dbx: insert query is nil")
	}
	rows := normalizedInsertRows(q)
	if err := validateInsertQuery(q, rows); err != nil {
		return BoundQuery{}, err
	}

	state := &renderState{dialect: d, args: make([]any, 0, len(rows)*4)}
	state.writeString(insertStatementPrefix(d, q))
	if err := renderInsertBody(state, q, rows); err != nil {
		return BoundQuery{}, err
	}
	if err := renderUpsert(state, q); err != nil {
		return BoundQuery{}, err
	}
	if err := renderReturning(state, q.ReturningItems); err != nil {
		return BoundQuery{}, err
	}
	if err := state.err(); err != nil {
		return BoundQuery{}, err
	}
	return state.BoundQuery(), nil
}

func validateInsertQuery(q *InsertQuery, rows [][]Assignment) error {
	switch {
	case q.Into.Name() == "":
		return errors.New("dbx: insert query requires target table")
	case len(rows) == 0 && q.Source == nil:
		return errors.New("dbx: insert query requires values or source query")
	case len(rows) > 0 && q.Source != nil:
		return errors.New("dbx: insert query cannot combine values and source query")
	case q.Source != nil && len(q.TargetColumns) == 0:
		return errors.New("dbx: insert-select requires target columns")
	default:
		return nil
	}
}

func insertStatementPrefix(d dialect.Dialect, q *InsertQuery) string {
	features := dialectFeatures(d)
	if features.InsertIgnoreForUpsertNothing && q.Upsert != nil && q.Upsert.DoNothing {
		return "INSERT IGNORE INTO "
	}
	return "INSERT INTO "
}

func renderInsertBody(state *renderState, q *InsertQuery, rows [][]Assignment) error {
	state.renderTable(q.Into)
	columns, err := resolveInsertColumns(q, rows)
	if err != nil {
		return err
	}
	if err := renderInsertColumns(state, columns); err != nil {
		return err
	}
	return renderInsertSourceOrValues(state, q, columns, rows)
}

func renderInsertColumns(state *renderState, columns []ColumnMeta) error {
	if len(columns) == 0 {
		return nil
	}
	state.writeString(" (")
	for i := range columns {
		if i > 0 {
			state.writeString(", ")
		}
		state.writeQuotedIdent(columns[i].Name)
	}
	state.writeByte(')')
	return nil
}

func renderInsertSourceOrValues(state *renderState, q *InsertQuery, columns []ColumnMeta, rows [][]Assignment) error {
	if q.Source != nil {
		state.writeByte(' ')
		return renderSelectQuery(state, q.Source)
	}
	return renderInsertValues(state, columns, rows)
}

func renderInsertValues(state *renderState, columns []ColumnMeta, rows [][]Assignment) error {
	orderedRows, err := orderInsertRows(columns, rows)
	if err != nil {
		return err
	}
	state.writeString(" VALUES ")
	for rowIndex, row := range orderedRows {
		if rowIndex > 0 {
			state.writeString(", ")
		}
		if err := renderInsertValueRow(state, row); err != nil {
			return err
		}
	}
	return nil
}

func renderInsertValueRow(state *renderState, row []Assignment) error {
	state.writeByte('(')
	for colIndex, assignment := range row {
		renderer, ok := assignment.(insertAssignmentRenderer)
		if !ok {
			return fmt.Errorf("dbx: unsupported insert assignment %T", assignment)
		}
		if colIndex > 0 {
			state.writeString(", ")
		}
		if err := renderer.renderAssignmentValue(state); err != nil {
			return err
		}
	}
	state.writeByte(')')
	return nil
}

func (q *UpdateQuery) Build(d dialect.Dialect) (BoundQuery, error) {
	if q == nil {
		return BoundQuery{}, errors.New("dbx: update query is nil")
	}
	if err := validateUpdateQuery(q); err != nil {
		return BoundQuery{}, err
	}

	state := &renderState{dialect: d, args: make([]any, 0, len(q.Assignments))}
	state.writeString("UPDATE ")
	state.renderTable(q.Table)
	state.writeString(" SET ")
	if err := renderUpdateAssignments(state, q.Assignments); err != nil {
		return BoundQuery{}, err
	}
	if err := renderOptionalWhere(state, q.WhereExp); err != nil {
		return BoundQuery{}, err
	}
	if err := renderReturning(state, q.ReturningItems); err != nil {
		return BoundQuery{}, err
	}
	if err := state.err(); err != nil {
		return BoundQuery{}, err
	}
	return state.BoundQuery(), nil
}

func validateUpdateQuery(q *UpdateQuery) error {
	switch {
	case q.Table.Name() == "":
		return errors.New("dbx: update query requires target table")
	case len(q.Assignments) == 0:
		return errors.New("dbx: update query requires assignments")
	default:
		return nil
	}
}

func renderUpdateAssignments(state *renderState, assignments []Assignment) error {
	for i, assignment := range assignments {
		if i > 0 {
			state.writeString(", ")
		}
		if err := renderAssignment(state, assignment); err != nil {
			return err
		}
	}
	return nil
}

func renderOptionalWhere(state *renderState, predicate Predicate) error {
	if predicate == nil {
		return nil
	}
	state.writeString(" WHERE ")
	return renderPredicate(state, predicate)
}

func (q *DeleteQuery) Build(d dialect.Dialect) (BoundQuery, error) {
	if q == nil {
		return BoundQuery{}, errors.New("dbx: delete query is nil")
	}
	if q.From.Name() == "" {
		return BoundQuery{}, errors.New("dbx: delete query requires target table")
	}

	state := &renderState{dialect: d, args: make([]any, 0, 4)}
	state.writeString("DELETE FROM ")
	state.renderTable(q.From)
	if err := renderOptionalWhere(state, q.WhereExp); err != nil {
		return BoundQuery{}, err
	}
	if err := renderReturning(state, q.ReturningItems); err != nil {
		return BoundQuery{}, err
	}
	if err := state.err(); err != nil {
		return BoundQuery{}, err
	}
	return state.BoundQuery(), nil
}

func normalizedInsertRows(q *InsertQuery) [][]Assignment {
	if len(q.Rows) > 0 {
		return q.Rows
	}
	if len(q.Assignments) > 0 {
		return [][]Assignment{q.Assignments}
	}
	return nil
}

func resolveInsertColumns(q *InsertQuery, rows [][]Assignment) ([]ColumnMeta, error) {
	if len(q.TargetColumns) > 0 {
		return resolveTargetColumns(q.TargetColumns)
	}
	if len(rows) == 0 {
		return nil, nil
	}
	return assignmentColumns(rows[0])
}

func assignmentColumns(assignments []Assignment) ([]ColumnMeta, error) {
	columns := collectionx.NewListWithCapacity[ColumnMeta](len(assignments))
	for _, assignment := range assignments {
		renderer, ok := assignment.(insertAssignmentRenderer)
		if !ok {
			return nil, fmt.Errorf("dbx: unsupported insert assignment %T", assignment)
		}
		columns.Add(renderer.assignmentColumn())
	}
	return columns.Values(), nil
}

func resolveTargetColumns(expressions []Expression) ([]ColumnMeta, error) {
	columns := collectionx.NewListWithCapacity[ColumnMeta](len(expressions))
	for _, expression := range expressions {
		column, ok := expression.(columnAccessor)
		if !ok {
			return nil, fmt.Errorf("dbx: unsupported target column expression %T", expression)
		}
		columns.Add(column.columnRef())
	}
	return columns.Values(), nil
}

func orderInsertRows(columns []ColumnMeta, rows [][]Assignment) ([][]Assignment, error) {
	orderedRows := collectionx.NewListWithCapacity[[]Assignment](len(rows))
	for _, row := range rows {
		orderedRow, err := orderInsertRow(columns, row)
		if err != nil {
			return nil, err
		}
		orderedRows.Add(orderedRow)
	}
	return orderedRows.Values(), nil
}

func orderInsertRow(columns []ColumnMeta, row []Assignment) ([]Assignment, error) {
	assignmentsByColumn := collectionx.NewMapWithCapacity[string, Assignment](len(row))
	for _, assignment := range row {
		renderer, ok := assignment.(insertAssignmentRenderer)
		if !ok {
			return nil, fmt.Errorf("dbx: unsupported insert assignment %T", assignment)
		}
		assignmentsByColumn.Set(renderer.assignmentColumn().Name, assignment)
	}

	orderedRow := collectionx.NewListWithCapacity[Assignment](len(columns))
	for i := range columns {
		assignment, ok := assignmentsByColumn.Get(columns[i].Name)
		if !ok {
			return nil, fmt.Errorf("dbx: missing value for insert column %s", columns[i].Name)
		}
		orderedRow.Add(assignment)
	}
	return orderedRow.Values(), nil
}
