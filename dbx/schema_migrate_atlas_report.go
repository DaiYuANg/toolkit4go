package dbx

import (
	"slices"

	atlasschema "ariga.io/atlas/sql/schema"
	"github.com/DaiYuANg/arcgo/collectionx"
)

func atlasReportFromChanges(changes []atlasschema.Change, compiled *atlasCompiledSchema, current *atlasschema.Schema) ValidationReport {
	diffs := atlasReportDiffMap(compiled.order)
	currentTables := atlasCurrentTablesByName(current)
	for _, change := range changes {
		atlasApplyChangeToReport(diffs, compiled, currentTables, change)
	}
	return atlasValidationReport(diffs)
}

func atlasReportDiffMap(order []string) collectionx.OrderedMap[string, *TableDiff] {
	diffs := collectionx.NewOrderedMap[string, *TableDiff]()
	for _, name := range order {
		diffs.Set(name, &TableDiff{Table: name})
	}
	return diffs
}

func atlasCurrentTablesByName(current *atlasschema.Schema) collectionx.Map[string, *atlasschema.Table] {
	currentTables := collectionx.NewMap[string, *atlasschema.Table]()
	if current == nil {
		return currentTables
	}
	for _, table := range current.Tables {
		currentTables.Set(table.Name, table)
	}
	return currentTables
}

func atlasApplyChangeToReport(diffs collectionx.OrderedMap[string, *TableDiff], compiled *atlasCompiledSchema, currentTables collectionx.Map[string, *atlasschema.Table], change atlasschema.Change) {
	switch c := change.(type) {
	case *atlasschema.AddTable:
		atlasApplyAddTableChange(diffs, compiled, c)
	case *atlasschema.ModifyTable:
		atlasApplyModifyTableChange(diffs, compiled, currentTables, c)
	}
}

func atlasApplyAddTableChange(diffs collectionx.OrderedMap[string, *TableDiff], compiled *atlasCompiledSchema, change *atlasschema.AddTable) {
	compiledTable, ok := compiled.tables.Get(change.T.Name)
	if !ok {
		return
	}
	diff, _ := diffs.Get(change.T.Name)
	diff.MissingTable = true
	diff.MissingColumns = slices.Clone(compiledTable.spec.Columns)
	diff.MissingIndexes = slices.Clone(compiledTable.spec.Indexes)
	diff.MissingForeignKeys = slices.Clone(compiledTable.spec.ForeignKeys)
	diff.MissingChecks = slices.Clone(compiledTable.spec.Checks)
	if compiledTable.spec.PrimaryKey != nil {
		diff.PrimaryKeyDiff = &PrimaryKeyDiff{
			Expected: new(clonePrimaryKeyMeta(*compiledTable.spec.PrimaryKey)),
			Issues:   []string{"table does not exist"},
		}
	}
}

func atlasApplyModifyTableChange(diffs collectionx.OrderedMap[string, *TableDiff], compiled *atlasCompiledSchema, currentTables collectionx.Map[string, *atlasschema.Table], change *atlasschema.ModifyTable) {
	compiledTable, ok := compiled.tables.Get(change.T.Name)
	if !ok {
		return
	}
	diff, _ := diffs.Get(change.T.Name)
	currentTable, _ := currentTables.Get(change.T.Name)
	for _, tableChange := range change.Changes {
		atlasApplyTableChangeToDiff(diff, compiledTable, currentTable, tableChange)
	}
}

func atlasValidationReport(diffs collectionx.OrderedMap[string, *TableDiff]) ValidationReport {
	reportTables := collectionx.NewListWithCapacity[TableDiff](diffs.Len())
	diffs.Range(func(_ string, diff *TableDiff) bool {
		reportTables.Add(*diff)
		return true
	})
	return ValidationReport{
		Tables:   reportTables.Values(),
		Backend:  ValidationBackendAtlas,
		Complete: true,
	}
}
