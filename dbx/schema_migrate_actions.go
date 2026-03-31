package dbx

import (
	"slices"
	"strings"
)

func mappedMigrationActions[T any](items []T, mapper func(T) MigrationAction) []MigrationAction {
	actions := make([]MigrationAction, 0, len(items))
	for i := range items {
		actions = append(actions, mapper(items[i]))
	}
	return actions
}

func columnDiffManualActions(diff TableDiff) []MigrationAction {
	return mappedMigrationActions(diff.ColumnDiffs, func(cd ColumnDiff) MigrationAction {
		return MigrationAction{
			Kind:    MigrationActionManual,
			Table:   diff.Table,
			Summary: "manual column migration required for " + cd.Column.Name + ": " + strings.Join(cd.Issues, "; "),
		}
	})
}

func columnDiffIssues(schemaDialect SchemaDialect, expected ColumnMeta, actual ColumnState) []string {
	issues := make([]string, 0, 4)
	expectedType := normalizeExpectedType(schemaDialect, expected)
	actualType := schemaDialect.NormalizeType(actual.Type)
	if expectedType != "" && actualType != "" && expectedType != actualType {
		issues = append(issues, "type mismatch: expected "+expectedType+" got "+actualType)
	}
	if !actual.PrimaryKey && expected.Nullable != actual.Nullable {
		issues = append(issues, "nullable mismatch")
	}
	if expected.AutoIncrement != actual.AutoIncrement {
		issues = append(issues, "auto increment mismatch")
	}
	if expected.DefaultValue != "" && normalizeDefault(expected.DefaultValue) != normalizeDefault(actual.DefaultValue) {
		issues = append(issues, "default mismatch")
	}
	return issues
}

func buildCreateTableAction(schemaDialect SchemaDialect, spec TableSpec) MigrationAction {
	statement, err := schemaDialect.BuildCreateTable(spec)
	if err != nil {
		return MigrationAction{
			Kind:    MigrationActionManual,
			Table:   spec.Name,
			Summary: "manual create table migration required: " + err.Error(),
		}
	}
	return MigrationAction{
		Kind:       MigrationActionCreateTable,
		Table:      spec.Name,
		Summary:    "create table " + spec.Name,
		Statement:  statement,
		Executable: true,
	}
}

func buildAddColumnAction(schemaDialect SchemaDialect, table string, column ColumnMeta) MigrationAction {
	statement, err := schemaDialect.BuildAddColumn(table, column)
	if err != nil {
		return MigrationAction{
			Kind:    MigrationActionManual,
			Table:   table,
			Summary: "manual add column migration required for " + column.Name + ": " + err.Error(),
		}
	}
	return MigrationAction{
		Kind:       MigrationActionAddColumn,
		Table:      table,
		Summary:    "add column " + column.Name,
		Statement:  statement,
		Executable: true,
	}
}

func buildCreateIndexAction(schemaDialect SchemaDialect, index IndexMeta) MigrationAction {
	statement, err := schemaDialect.BuildCreateIndex(index)
	if err != nil {
		return MigrationAction{
			Kind:    MigrationActionManual,
			Table:   index.Table,
			Summary: "manual create index migration required for " + index.Name + ": " + err.Error(),
		}
	}
	return MigrationAction{
		Kind:       MigrationActionCreateIndex,
		Table:      index.Table,
		Summary:    "create index " + index.Name,
		Statement:  statement,
		Executable: true,
	}
}

func buildAddForeignKeyAction(schemaDialect SchemaDialect, table string, foreignKey ForeignKeyMeta) MigrationAction {
	statement, err := schemaDialect.BuildAddForeignKey(table, foreignKey)
	if err != nil {
		return MigrationAction{
			Kind:    MigrationActionManual,
			Table:   table,
			Summary: "manual add foreign key migration required for " + foreignKey.Name + ": " + err.Error(),
		}
	}
	return MigrationAction{
		Kind:       MigrationActionAddForeignKey,
		Table:      table,
		Summary:    "add foreign key " + foreignKey.Name,
		Statement:  statement,
		Executable: true,
	}
}

func buildAddCheckAction(schemaDialect SchemaDialect, table string, check CheckMeta) MigrationAction {
	statement, err := schemaDialect.BuildAddCheck(table, check)
	if err != nil {
		return MigrationAction{
			Kind:    MigrationActionManual,
			Table:   table,
			Summary: "manual add check migration required for " + check.Name + ": " + err.Error(),
		}
	}
	return MigrationAction{
		Kind:       MigrationActionAddCheck,
		Table:      table,
		Summary:    "add check " + check.Name,
		Statement:  statement,
		Executable: true,
	}
}

func primaryKeyIssues(expected *PrimaryKeyMeta, actual *PrimaryKeyState) []string {
	switch {
	case expected == nil && actual == nil:
		return nil
	case expected == nil:
		return []string{"unexpected primary key present"}
	case actual == nil:
		return []string{"missing primary key"}
	case columnsKey(expected.Columns) != columnsKey(actual.Columns):
		return []string{"primary key columns mismatch"}
	default:
		return nil
	}
}

func indexKey(unique bool, columns []string) string {
	prefix := "idx:"
	if unique {
		prefix = "ux:"
	}
	return prefix + columnsKey(columns)
}

func foreignKeyKey(meta ForeignKeyMeta) string {
	return columnsKey(meta.Columns) + "->" + meta.TargetTable + ":" + columnsKey(meta.TargetColumns) + ":" + string(normalizeReferentialAction(meta.OnDelete)) + ":" + string(normalizeReferentialAction(meta.OnUpdate))
}

func foreignKeyKeyFromState(state ForeignKeyState) string {
	return columnsKey(state.Columns) + "->" + state.TargetTable + ":" + columnsKey(state.TargetColumns) + ":" + string(normalizeReferentialAction(state.OnDelete)) + ":" + string(normalizeReferentialAction(state.OnUpdate))
}

func checkKey(expression string) string {
	return normalizeCheckExpression(expression)
}

func columnsKey(columns []string) string {
	return strings.Join(columns, ",")
}

func normalizeReferentialAction(action ReferentialAction) ReferentialAction {
	if strings.TrimSpace(string(action)) == "" {
		return ReferentialNoAction
	}
	return action
}

func clonePrimaryKeyState(state PrimaryKeyState) PrimaryKeyState {
	state.Columns = slices.Clone(state.Columns)
	return state
}
