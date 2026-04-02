package dbx

import (
	atlasschema "ariga.io/atlas/sql/schema"

	"github.com/DaiYuANg/arcgo/collectionx"
)

type atlasTableChangeHandler func(*TableDiff, *atlasCompiledTable, *atlasschema.Table, atlasschema.Change) bool

var atlasTableChangeHandlers = []atlasTableChangeHandler{
	handleAtlasAddColumnChange,
	handleAtlasModifyColumnChange,
	handleAtlasRenameColumnChange,
	handleAtlasDropColumnChange,
	handleAtlasAddIndexChange,
	handleAtlasModifyIndexChange,
	handleAtlasRenameIndexChange,
	handleAtlasDropIndexChange,
	handleAtlasAddForeignKeyChange,
	handleAtlasModifyForeignKeyChange,
	handleAtlasDropForeignKeyChange,
	handleAtlasAddCheckChange,
	handleAtlasModifyCheckChange,
	handleAtlasDropCheckChange,
	handleAtlasAddPrimaryKeyChange,
	handleAtlasModifyOrDropPrimaryKeyChange,
}

func atlasApplyTableChangeToDiff(diff *TableDiff, compiled *atlasCompiledTable, current *atlasschema.Table, change atlasschema.Change) {
	for _, handler := range atlasTableChangeHandlers {
		if handler(diff, compiled, current, change) {
			return
		}
	}
}

func handleAtlasAddColumnChange(diff *TableDiff, compiled *atlasCompiledTable, _ *atlasschema.Table, change atlasschema.Change) bool {
	current, ok := change.(*atlasschema.AddColumn)
	if !ok {
		return false
	}
	atlasHandleAddColumn(diff, compiled, current)
	return true
}

func handleAtlasModifyColumnChange(diff *TableDiff, compiled *atlasCompiledTable, _ *atlasschema.Table, change atlasschema.Change) bool {
	current, ok := change.(*atlasschema.ModifyColumn)
	if !ok {
		return false
	}
	atlasHandleModifyColumn(diff, compiled, current)
	return true
}

func handleAtlasRenameColumnChange(diff *TableDiff, _ *atlasCompiledTable, _ *atlasschema.Table, change atlasschema.Change) bool {
	current, ok := change.(*atlasschema.RenameColumn)
	if !ok {
		return false
	}
	atlasHandleRenameColumn(diff, current)
	return true
}

func handleAtlasDropColumnChange(diff *TableDiff, _ *atlasCompiledTable, _ *atlasschema.Table, change atlasschema.Change) bool {
	current, ok := change.(*atlasschema.DropColumn)
	if !ok {
		return false
	}
	atlasHandleDropColumn(diff, current)
	return true
}

func handleAtlasAddIndexChange(diff *TableDiff, compiled *atlasCompiledTable, _ *atlasschema.Table, change atlasschema.Change) bool {
	current, ok := change.(*atlasschema.AddIndex)
	if !ok {
		return false
	}
	atlasHandleAddIndex(diff, compiled, current)
	return true
}

func handleAtlasModifyIndexChange(diff *TableDiff, compiled *atlasCompiledTable, _ *atlasschema.Table, change atlasschema.Change) bool {
	current, ok := change.(*atlasschema.ModifyIndex)
	if !ok {
		return false
	}
	atlasHandleModifyIndex(diff, compiled, current)
	return true
}

func handleAtlasRenameIndexChange(diff *TableDiff, _ *atlasCompiledTable, _ *atlasschema.Table, change atlasschema.Change) bool {
	current, ok := change.(*atlasschema.RenameIndex)
	if !ok {
		return false
	}
	atlasHandleRenameIndex(diff, current)
	return true
}

func handleAtlasDropIndexChange(diff *TableDiff, _ *atlasCompiledTable, _ *atlasschema.Table, change atlasschema.Change) bool {
	current, ok := change.(*atlasschema.DropIndex)
	if !ok {
		return false
	}
	atlasHandleDropIndex(diff, current)
	return true
}

func handleAtlasAddForeignKeyChange(diff *TableDiff, compiled *atlasCompiledTable, _ *atlasschema.Table, change atlasschema.Change) bool {
	current, ok := change.(*atlasschema.AddForeignKey)
	if !ok {
		return false
	}
	atlasHandleAddForeignKey(diff, compiled, current)
	return true
}

func handleAtlasModifyForeignKeyChange(diff *TableDiff, compiled *atlasCompiledTable, _ *atlasschema.Table, change atlasschema.Change) bool {
	current, ok := change.(*atlasschema.ModifyForeignKey)
	if !ok {
		return false
	}
	atlasHandleModifyForeignKey(diff, compiled, current)
	return true
}

func handleAtlasDropForeignKeyChange(diff *TableDiff, _ *atlasCompiledTable, _ *atlasschema.Table, change atlasschema.Change) bool {
	current, ok := change.(*atlasschema.DropForeignKey)
	if !ok {
		return false
	}
	atlasHandleDropForeignKey(diff, current)
	return true
}

func handleAtlasAddCheckChange(diff *TableDiff, compiled *atlasCompiledTable, _ *atlasschema.Table, change atlasschema.Change) bool {
	current, ok := change.(*atlasschema.AddCheck)
	if !ok {
		return false
	}
	atlasHandleAddCheck(diff, compiled, current)
	return true
}

func handleAtlasModifyCheckChange(diff *TableDiff, compiled *atlasCompiledTable, _ *atlasschema.Table, change atlasschema.Change) bool {
	current, ok := change.(*atlasschema.ModifyCheck)
	if !ok {
		return false
	}
	atlasHandleModifyCheck(diff, compiled, current)
	return true
}

func handleAtlasDropCheckChange(diff *TableDiff, _ *atlasCompiledTable, _ *atlasschema.Table, change atlasschema.Change) bool {
	current, ok := change.(*atlasschema.DropCheck)
	if !ok {
		return false
	}
	atlasHandleDropCheck(diff, current)
	return true
}

func handleAtlasAddPrimaryKeyChange(diff *TableDiff, compiled *atlasCompiledTable, _ *atlasschema.Table, change atlasschema.Change) bool {
	_, ok := change.(*atlasschema.AddPrimaryKey)
	if !ok {
		return false
	}
	atlasHandleAddPrimaryKey(diff, compiled)
	return true
}

func handleAtlasModifyOrDropPrimaryKeyChange(diff *TableDiff, compiled *atlasCompiledTable, current *atlasschema.Table, change atlasschema.Change) bool {
	switch change.(type) {
	case *atlasschema.ModifyPrimaryKey, *atlasschema.DropPrimaryKey:
		atlasHandleModifyOrDropPrimaryKey(diff, compiled, current)
		return true
	default:
		return false
	}
}

func atlasHandleAddColumn(diff *TableDiff, compiled *atlasCompiledTable, change *atlasschema.AddColumn) {
	if column, ok := compiled.columnsByName.Get(change.C.Name); ok {
		diff.MissingColumns.Add(column)
	}
}

func atlasHandleModifyColumn(diff *TableDiff, compiled *atlasCompiledTable, change *atlasschema.ModifyColumn) {
	name := change.To.Name
	if name == "" {
		name = change.From.Name
	}
	column, ok := compiled.columnsByName.Get(name)
	if !ok {
		column = ColumnMeta{Name: name, Table: diff.Table}
	}
	diff.ColumnDiffs.Add(ColumnDiff{Column: column, Issues: collectionx.NewList(atlasColumnChangeIssue(change.Change))})
}

func atlasHandleRenameColumn(diff *TableDiff, change *atlasschema.RenameColumn) {
	diff.ColumnDiffs.Add(ColumnDiff{Column: ColumnMeta{Name: change.To.Name, Table: diff.Table}, Issues: collectionx.NewList("manual column rename migration required")})
}

func atlasHandleDropColumn(diff *TableDiff, change *atlasschema.DropColumn) {
	diff.ColumnDiffs.Add(ColumnDiff{Column: ColumnMeta{Name: change.C.Name, Table: diff.Table}, Issues: collectionx.NewList("manual column removal migration required")})
}

func atlasHandleAddIndex(diff *TableDiff, compiled *atlasCompiledTable, change *atlasschema.AddIndex) {
	if index, ok := atlasFindIndexMeta(compiled, change.I); ok {
		diff.MissingIndexes.Add(index)
	}
}

func atlasHandleModifyIndex(diff *TableDiff, compiled *atlasCompiledTable, change *atlasschema.ModifyIndex) {
	if index, ok := atlasFindIndexMeta(compiled, change.To); ok {
		diff.MissingIndexes.Add(index)
		return
	}
	diff.ColumnDiffs.Add(ColumnDiff{Column: ColumnMeta{Name: change.To.Name, Table: diff.Table}, Issues: collectionx.NewList("manual index modification required")})
}

func atlasHandleRenameIndex(diff *TableDiff, change *atlasschema.RenameIndex) {
	diff.ColumnDiffs.Add(ColumnDiff{Column: ColumnMeta{Name: change.To.Name, Table: diff.Table}, Issues: collectionx.NewList("manual index rename migration required")})
}

func atlasHandleDropIndex(diff *TableDiff, change *atlasschema.DropIndex) {
	diff.ColumnDiffs.Add(ColumnDiff{Column: ColumnMeta{Name: change.I.Name, Table: diff.Table}, Issues: collectionx.NewList("manual index removal migration required")})
}

func atlasHandleAddForeignKey(diff *TableDiff, compiled *atlasCompiledTable, change *atlasschema.AddForeignKey) {
	if foreignKey, ok := atlasFindForeignKeyMeta(compiled, change.F); ok {
		diff.MissingForeignKeys.Add(foreignKey)
	}
}

func atlasHandleModifyForeignKey(diff *TableDiff, compiled *atlasCompiledTable, change *atlasschema.ModifyForeignKey) {
	if foreignKey, ok := atlasFindForeignKeyMeta(compiled, change.To); ok {
		diff.MissingForeignKeys.Add(foreignKey)
	}
}

func atlasHandleDropForeignKey(diff *TableDiff, change *atlasschema.DropForeignKey) {
	diff.ColumnDiffs.Add(ColumnDiff{Column: ColumnMeta{Name: change.F.Symbol, Table: diff.Table}, Issues: collectionx.NewList("manual foreign key removal migration required")})
}

func atlasHandleAddCheck(diff *TableDiff, compiled *atlasCompiledTable, change *atlasschema.AddCheck) {
	if check, ok := atlasFindCheckMeta(compiled, change.C); ok {
		diff.MissingChecks.Add(check)
	}
}

func atlasHandleModifyCheck(diff *TableDiff, compiled *atlasCompiledTable, change *atlasschema.ModifyCheck) {
	if check, ok := atlasFindCheckMeta(compiled, change.To); ok {
		diff.MissingChecks.Add(check)
	}
}

func atlasHandleDropCheck(diff *TableDiff, change *atlasschema.DropCheck) {
	diff.ColumnDiffs.Add(ColumnDiff{Column: ColumnMeta{Name: change.C.Name, Table: diff.Table}, Issues: collectionx.NewList("manual check removal migration required")})
}

func atlasHandleAddPrimaryKey(diff *TableDiff, compiled *atlasCompiledTable) {
	diff.PrimaryKeyDiff = &PrimaryKeyDiff{
		Expected: compiled.spec.PrimaryKey,
		Actual:   atlasPrimaryKeyState(nil),
		Issues:   collectionx.NewList("missing primary key"),
	}
}

func atlasHandleModifyOrDropPrimaryKey(diff *TableDiff, compiled *atlasCompiledTable, current *atlasschema.Table) {
	var actual *PrimaryKeyState
	if current != nil {
		actual = atlasPrimaryKeyState(current)
	}
	var expected *PrimaryKeyMeta
	if compiled.spec.PrimaryKey != nil {
		expected = new(clonePrimaryKeyMeta(*compiled.spec.PrimaryKey))
	}
	diff.PrimaryKeyDiff = &PrimaryKeyDiff{Expected: expected, Actual: actual, Issues: collectionx.NewList("primary key migration required")}
}
