package dbx_test

import (
	"strings"
)

func toColumnState(column ColumnMeta) ColumnState {
	typeName := column.SQLType
	if typeName == "" {
		typeName = InferTypeNameForTest(column)
	}
	return ColumnState{
		Name:          column.Name,
		Type:          strings.ToLower(typeName),
		Nullable:      column.Nullable,
		PrimaryKey:    column.PrimaryKey,
		AutoIncrement: column.AutoIncrement,
		DefaultValue:  column.DefaultValue,
	}
}

func toIndexStates(indexes []IndexMeta) []IndexState {
	items := make([]IndexState, len(indexes))
	for i, index := range indexes {
		items[i] = IndexState{
			Name:    index.Name,
			Columns: append([]string(nil), index.Columns...),
			Unique:  index.Unique,
		}
	}
	return items
}

func toForeignKeyStates(foreignKeys []ForeignKeyMeta) []ForeignKeyState {
	items := make([]ForeignKeyState, len(foreignKeys))
	for i := range foreignKeys {
		foreignKey := &foreignKeys[i]
		items[i] = ForeignKeyState{
			Name:          foreignKey.Name,
			Columns:       append([]string(nil), foreignKey.Columns...),
			TargetTable:   foreignKey.TargetTable,
			TargetColumns: append([]string(nil), foreignKey.TargetColumns...),
			OnDelete:      foreignKey.OnDelete,
			OnUpdate:      foreignKey.OnUpdate,
		}
	}
	return items
}

func toCheckStates(checks []CheckMeta) []CheckState {
	items := make([]CheckState, len(checks))
	for i, check := range checks {
		items[i] = CheckState{Name: check.Name, Expression: check.Expression}
	}
	return items
}
