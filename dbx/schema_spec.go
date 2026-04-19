package dbx

import (
	schemax "github.com/DaiYuANg/arcgo/dbx/schema"
	"strings"

	"github.com/DaiYuANg/arcgo/collectionx"
)

func buildTableSpec(def schemaDefinition) schemax.TableSpec {
	indexes := deriveIndexes(def)
	foreignKeys := deriveForeignKeys(def)
	checks := deriveChecks(def)
	return schemax.TableSpec{
		Name:        def.table.Name(),
		Columns:     cloneColumnMetas(def.columns),
		Indexes:     collectionx.NewListWithCapacity(len(indexes), indexes...),
		PrimaryKey:  derivePrimaryKey(def),
		ForeignKeys: collectionx.NewListWithCapacity(len(foreignKeys), foreignKeys...),
		Checks:      collectionx.NewListWithCapacity(len(checks), checks...),
	}
}

func deriveIndexes(def schemaDefinition) []schemax.IndexMeta {
	indexes := collectionx.NewOrderedMap[string, schemax.IndexMeta]()
	def.indexes.Range(func(_ int, index schemax.IndexMeta) bool {
		indexes.Set(indexKey(index.Unique, index.Columns), cloneIndexMeta(index))
		return true
	})
	deriveColumnIndexes(def, indexes)
	items := make([]schemax.IndexMeta, 0, indexes.Len())
	indexes.Range(func(_ string, value schemax.IndexMeta) bool {
		items = append(items, cloneIndexMeta(value))
		return true
	})
	return items
}

func deriveColumnIndexes(def schemaDefinition, indexes collectionx.OrderedMap[string, schemax.IndexMeta]) {
	tableName := def.table.Name()
	def.columns.Range(func(_ int, column schemax.ColumnMeta) bool {
		if !shouldDeriveColumnIndex(column) {
			return true
		}
		meta := schemax.IndexMeta{
			Name:    indexNameForColumn(tableName, column),
			Table:   tableName,
			Columns: collectionx.NewList(column.Name),
			Unique:  column.Unique,
		}
		indexes.Set(indexKey(meta.Unique, meta.Columns), meta)
		return true
	})
}

func shouldDeriveColumnIndex(column schemax.ColumnMeta) bool {
	return !column.PrimaryKey && (column.Unique || column.Indexed)
}

func indexNameForColumn(table string, column schemax.ColumnMeta) string {
	prefix := "idx"
	if column.Unique {
		prefix = "ux"
	}
	return prefix + "_" + table + "_" + column.Name
}

func derivePrimaryKey(def schemaDefinition) *schemax.PrimaryKeyMeta {
	tableName := def.table.Name()
	if def.primaryKey != nil {
		copyPrimary := clonePrimaryKeyMeta(*def.primaryKey)
		if copyPrimary.Name == "" {
			copyPrimary.Name = "pk_" + tableName
		}
		if copyPrimary.Table == "" {
			copyPrimary.Table = tableName
		}
		return &copyPrimary
	}

	columns := collectionx.FilterMapList(def.columns, func(_ int, column schemax.ColumnMeta) (string, bool) {
		return column.Name, column.PrimaryKey
	})
	if columns.Len() == 0 {
		return nil
	}
	return &schemax.PrimaryKeyMeta{
		Name:    "pk_" + tableName,
		Table:   tableName,
		Columns: columns,
	}
}

func deriveForeignKeys(def schemaDefinition) []schemax.ForeignKeyMeta {
	foreignKeys := collectionx.NewOrderedMap[string, schemax.ForeignKeyMeta]()
	explicitColumns := collectionx.NewSet[string]()
	deriveExplicitForeignKeys(def, foreignKeys, explicitColumns)
	deriveRelationForeignKeys(def, foreignKeys, explicitColumns)
	items := make([]schemax.ForeignKeyMeta, 0, foreignKeys.Len())
	foreignKeys.Range(func(_ string, value schemax.ForeignKeyMeta) bool {
		items = append(items, cloneForeignKeyMeta(value))
		return true
	})
	return items
}

func deriveExplicitForeignKeys(def schemaDefinition, foreignKeys collectionx.OrderedMap[string, schemax.ForeignKeyMeta], explicitColumns collectionx.Set[string]) {
	tableName := def.table.Name()
	def.columns.Range(func(_ int, column schemax.ColumnMeta) bool {
		if column.References == nil {
			return true
		}
		explicitColumns.Add(column.Name)
		meta := schemax.ForeignKeyMeta{
			Name:          "fk_" + tableName + "_" + column.Name,
			Table:         tableName,
			Columns:       collectionx.NewList(column.Name),
			TargetTable:   column.References.TargetTable,
			TargetColumns: collectionx.NewList(column.References.TargetColumn),
			OnDelete:      column.References.OnDelete,
			OnUpdate:      column.References.OnUpdate,
		}
		foreignKeys.Set(foreignKeyKey(meta), meta)
		return true
	})
}

func deriveRelationForeignKeys(def schemaDefinition, foreignKeys collectionx.OrderedMap[string, schemax.ForeignKeyMeta], explicitColumns collectionx.Set[string]) {
	tableName := def.table.Name()
	def.relations.Range(func(_ int, relation schemax.RelationMeta) bool {
		if !shouldDeriveRelationForeignKey(def, relation, explicitColumns) {
			return true
		}
		meta := schemax.ForeignKeyMeta{
			Name:          "fk_" + tableName + "_" + relation.LocalColumn,
			Table:         tableName,
			Columns:       collectionx.NewList(relation.LocalColumn),
			TargetTable:   relation.TargetTable,
			TargetColumns: collectionx.NewList(relation.TargetColumn),
		}
		key := foreignKeyKey(meta)
		if _, exists := foreignKeys.Get(key); !exists {
			foreignKeys.Set(key, meta)
		}
		return true
	})
}

func shouldDeriveRelationForeignKey(def schemaDefinition, relation schemax.RelationMeta, explicitColumns collectionx.Set[string]) bool {
	if relation.Kind != schemax.RelationBelongsTo {
		return false
	}
	if relation.LocalColumn == "" || relation.TargetColumn == "" || relation.TargetTable == "" {
		return false
	}
	if explicitColumns.Contains(relation.LocalColumn) {
		return false
	}
	_, ok := def.columnByName(relation.LocalColumn)
	return ok
}

func deriveChecks(def schemaDefinition) []schemax.CheckMeta {
	return collectionx.MapList(def.checks, func(_ int, check schemax.CheckMeta) schemax.CheckMeta {
		return cloneCheckMeta(check)
	}).Values()
}

func indexKey(unique bool, columns collectionx.List[string]) string {
	prefix := "idx:"
	if unique {
		prefix = "ux:"
	}
	return prefix + columnsKey(columns)
}

func foreignKeyKey(meta schemax.ForeignKeyMeta) string {
	return columnsKey(meta.Columns) + "->" + meta.TargetTable + ":" + columnsKey(meta.TargetColumns) + ":" + string(normalizeReferentialAction(meta.OnDelete)) + ":" + string(normalizeReferentialAction(meta.OnUpdate))
}

func columnsKey(columns collectionx.List[string]) string {
	return columns.Join(",")
}

func normalizeReferentialAction(action schemax.ReferentialAction) schemax.ReferentialAction {
	if strings.TrimSpace(string(action)) == "" {
		return schemax.ReferentialNoAction
	}
	return action
}
