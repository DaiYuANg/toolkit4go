package dbx

import (
	"context"
	"slices"

	"github.com/DaiYuANg/arcgo/collectionx"
	"github.com/samber/lo"
)

func diffSchema(ctx context.Context, schemaDialect SchemaDialect, session Session, schema SchemaResource) (TableDiff, error) {
	spec := buildTableSpec(schema.schemaRef())
	actual, err := schemaDialect.InspectTable(ctx, session, spec.Name)
	if err != nil {
		return TableDiff{}, wrapDBError("inspect schema table", err)
	}
	if !actual.Exists {
		return missingTableDiff(spec), nil
	}
	return existingTableDiff(schemaDialect, spec, actual), nil
}

func missingTableDiff(spec TableSpec) TableDiff {
	diff := newTableDiff(spec.Name)
	diff.MissingTable = true
	diff.MissingColumns = collectionx.NewListWithCapacity(len(spec.Columns), slices.Clone(spec.Columns)...)
	diff.MissingIndexes = collectionx.NewListWithCapacity(len(spec.Indexes), slices.Clone(spec.Indexes)...)
	diff.MissingForeignKeys = collectionx.NewListWithCapacity(len(spec.ForeignKeys), slices.Clone(spec.ForeignKeys)...)
	diff.MissingChecks = collectionx.NewListWithCapacity(len(spec.Checks), slices.Clone(spec.Checks)...)
	if spec.PrimaryKey != nil {
		diff.PrimaryKeyDiff = &PrimaryKeyDiff{
			Expected: new(clonePrimaryKeyMeta(*spec.PrimaryKey)),
			Issues:   collectionx.NewList("table does not exist"),
		}
	}
	return diff
}

func existingTableDiff(schemaDialect SchemaDialect, spec TableSpec, actual TableState) TableDiff {
	diff := newTableDiff(spec.Name)
	actualColumns := lo.SliceToMap(actual.Columns, func(column ColumnState) (string, ColumnState) {
		return column.Name, column
	})
	diffColumns(schemaDialect, spec.Columns, actualColumns, &diff)
	diffPrimaryKey(spec.PrimaryKey, actual.PrimaryKey, &diff)
	diffIndexes(spec.Indexes, actual.Indexes, &diff)
	diffForeignKeys(spec.ForeignKeys, actual.ForeignKeys, &diff)
	diffChecks(spec.Checks, actual.Checks, &diff)
	return diff
}

func diffColumns(schemaDialect SchemaDialect, expectedColumns []ColumnMeta, actualColumns map[string]ColumnState, diff *TableDiff) {
	missingColumns := collectionx.NewListWithCapacity[ColumnMeta](len(expectedColumns))
	columnDiffs := collectionx.NewListWithCapacity[ColumnDiff](len(expectedColumns))
	for i := range expectedColumns {
		expected := expectedColumns[i]
		actualColumn, ok := actualColumns[expected.Name]
		if !ok {
			missingColumns.Add(expected)
			continue
		}
		issues := columnDiffIssues(schemaDialect, expected, actualColumn)
		if len(issues) == 0 {
			continue
		}
		columnDiffs.Add(ColumnDiff{Column: expected, Issues: collectionx.NewListWithCapacity(len(issues), issues...)})
	}
	diff.MissingColumns = missingColumns
	diff.ColumnDiffs = columnDiffs
}

func diffPrimaryKey(expected *PrimaryKeyMeta, actual *PrimaryKeyState, diff *TableDiff) {
	issues := primaryKeyIssues(expected, actual)
	if len(issues) == 0 {
		return
	}
	diff.PrimaryKeyDiff = &PrimaryKeyDiff{
		Expected: clonePrimaryKeyMetaPtr(expected),
		Actual:   clonePrimaryKeyStatePtr(actual),
		Issues:   collectionx.NewListWithCapacity(len(issues), issues...),
	}
}

func clonePrimaryKeyMetaPtr(meta *PrimaryKeyMeta) *PrimaryKeyMeta {
	if meta == nil {
		return nil
	}
	return new(clonePrimaryKeyMeta(*meta))
}

func clonePrimaryKeyStatePtr(state *PrimaryKeyState) *PrimaryKeyState {
	if state == nil {
		return nil
	}
	return new(clonePrimaryKeyState(*state))
}

func diffIndexes(expected []IndexMeta, actual []IndexState, diff *TableDiff) {
	actualIndexes := lo.SliceToMap(actual, func(index IndexState) (string, IndexState) {
		return indexKey(index.Unique, index.Columns), index
	})
	diff.MissingIndexes = missingByKey(expected, actualIndexes, func(index IndexMeta) string {
		return indexKey(index.Unique, index.Columns)
	})
}

func diffForeignKeys(expected []ForeignKeyMeta, actual []ForeignKeyState, diff *TableDiff) {
	actualForeignKeys := lo.SliceToMap(actual, func(foreignKey ForeignKeyState) (string, ForeignKeyState) {
		return foreignKeyKeyFromState(foreignKey), foreignKey
	})
	diff.MissingForeignKeys = missingByKey(expected, actualForeignKeys, foreignKeyKey)
}

func diffChecks(expected []CheckMeta, actual []CheckState, diff *TableDiff) {
	actualChecks := lo.SliceToMap(actual, func(check CheckState) (string, CheckState) {
		return checkKey(check.Expression), check
	})
	diff.MissingChecks = missingByKey(expected, actualChecks, func(check CheckMeta) string {
		return checkKey(check.Expression)
	})
}

func missingByKey[T any, S any](expected []T, actual map[string]S, key func(T) string) collectionx.List[T] {
	return collectionx.NewList(lo.Filter(expected, func(item T, _ int) bool {
		_, ok := actual[key(item)]
		return !ok
	})...)
}

func buildTableSpec(def schemaDefinition) TableSpec {
	return TableSpec{
		Name:        def.table.name,
		Columns:     slices.Clone(def.columns),
		Indexes:     deriveIndexes(def),
		PrimaryKey:  derivePrimaryKey(def),
		ForeignKeys: deriveForeignKeys(def),
		Checks:      deriveChecks(def),
	}
}
