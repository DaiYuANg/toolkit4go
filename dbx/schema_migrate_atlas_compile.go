package dbx

import (
	"strings"

	atlasmigrate "ariga.io/atlas/sql/migrate"
	atlasschema "ariga.io/atlas/sql/schema"
	"github.com/DaiYuANg/arcgo/collectionx"
	"github.com/samber/lo"
)

func compileAtlasSchema(dialectName string, driver atlasmigrate.Driver, schemaName string, schemas []SchemaResource) *atlasCompiledSchema {
	compiled := newAtlasCompiledSchema(schemaName, schemas)
	for _, resource := range schemas {
		compileAtlasResource(compiled, dialectName, driver, resource)
	}
	attachCompiledTableConstraints(compiled)
	return compiled
}

func newAtlasCompiledSchema(schemaName string, schemas []SchemaResource) *atlasCompiledSchema {
	atlasSchema := atlasschema.New(schemaName)
	return &atlasCompiledSchema{
		schema:    atlasSchema,
		tables:    collectionx.NewMapWithCapacity[string, *atlasCompiledTable](len(schemas)),
		externals: collectionx.NewMap[string, *atlasschema.Table](),
		order: lo.Map(schemas, func(schema SchemaResource, _ int) string {
			return schema.tableRef().TableName()
		}),
	}
}

func compileAtlasResource(compiled *atlasCompiledSchema, dialectName string, driver atlasmigrate.Driver, resource SchemaResource) {
	spec := buildTableSpec(resource.schemaRef())
	table := atlasschema.NewTable(spec.Name).SetSchema(compiled.schema)
	compiledTable := newAtlasCompiledTable(spec, table)
	compileAtlasTableColumns(compiledTable, dialectName, driver)
	compileAtlasTableMetadata(compiledTable)
	compiled.schema.AddTables(table)
	compiled.tables.Set(spec.Name, compiledTable)
}

func newAtlasCompiledTable(spec TableSpec, table *atlasschema.Table) *atlasCompiledTable {
	return &atlasCompiledTable{
		spec:              spec,
		table:             table,
		columnsByName:     collectionx.NewMapWithCapacity[string, ColumnMeta](len(spec.Columns)),
		indexesByName:     collectionx.NewMapWithCapacity[string, IndexMeta](len(spec.Indexes)),
		indexesByKey:      collectionx.NewMapWithCapacity[string, IndexMeta](len(spec.Indexes)),
		foreignKeysByName: collectionx.NewMapWithCapacity[string, ForeignKeyMeta](len(spec.ForeignKeys)),
		foreignKeysByKey:  collectionx.NewMapWithCapacity[string, ForeignKeyMeta](len(spec.ForeignKeys)),
		checksByName:      collectionx.NewMapWithCapacity[string, CheckMeta](len(spec.Checks)),
		checksByExpr:      collectionx.NewMapWithCapacity[string, CheckMeta](len(spec.Checks)),
	}
}

func compileAtlasTableColumns(compiledTable *atlasCompiledTable, dialectName string, driver atlasmigrate.Driver) {
	for i := range compiledTable.spec.Columns {
		column := compiledTable.spec.Columns[i]
		atlasColumn := compileAtlasColumn(dialectName, driver, column)
		compiledTable.table.AddColumns(atlasColumn)
		compiledTable.columnsByName.Set(column.Name, column)
	}
}

func compileAtlasTableMetadata(compiledTable *atlasCompiledTable) {
	for _, index := range compiledTable.spec.Indexes {
		compiledTable.indexesByName.Set(index.Name, index)
		compiledTable.indexesByKey.Set(indexKey(index.Unique, index.Columns), index)
	}
	for i := range compiledTable.spec.ForeignKeys {
		foreignKey := compiledTable.spec.ForeignKeys[i]
		compiledTable.foreignKeysByName.Set(foreignKey.Name, foreignKey)
		compiledTable.foreignKeysByKey.Set(foreignKeyKey(foreignKey), foreignKey)
	}
	for _, check := range compiledTable.spec.Checks {
		compiledTable.checksByName.Set(check.Name, check)
		compiledTable.checksByExpr.Set(checkKey(check.Expression), check)
	}
}

func attachCompiledTableConstraints(compiled *atlasCompiledSchema) {
	compiled.tables.Range(func(_ string, table *atlasCompiledTable) bool {
		attachCompiledPrimaryKey(table)
		attachCompiledIndexes(table)
		attachCompiledForeignKeys(compiled, table)
		attachCompiledChecks(table)
		return true
	})
}

func attachCompiledPrimaryKey(table *atlasCompiledTable) {
	if table.spec.PrimaryKey == nil {
		return
	}
	if primaryKey := atlasPrimaryKeyForSpec(table.table, *table.spec.PrimaryKey); primaryKey != nil {
		table.table.SetPrimaryKey(primaryKey)
	}
}

func attachCompiledIndexes(table *atlasCompiledTable) {
	for _, index := range table.spec.Indexes {
		if atlasIndex := atlasIndexForSpec(table.table, index); atlasIndex != nil {
			table.table.AddIndexes(atlasIndex)
		}
	}
}

func attachCompiledForeignKeys(compiled *atlasCompiledSchema, table *atlasCompiledTable) {
	for i := range table.spec.ForeignKeys {
		foreignKey := table.spec.ForeignKeys[i]
		if atlasForeignKey := compiled.atlasForeignKeyForSpec(table.table, foreignKey); atlasForeignKey != nil {
			table.table.AddForeignKeys(atlasForeignKey)
		}
	}
}

func attachCompiledChecks(table *atlasCompiledTable) {
	for _, check := range table.spec.Checks {
		table.table.AddChecks(atlasschema.NewCheck().SetName(check.Name).SetExpr(check.Expression))
	}
}

func compileAtlasColumn(dialectName string, driver atlasmigrate.Driver, column ColumnMeta) *atlasschema.Column {
	rawType := atlasColumnRawType(dialectName, column)
	atlasColumn := atlasschema.NewColumn(column.Name)
	atlasColumn.Type = &atlasschema.ColumnType{
		Type: atlasColumnType(driver, rawType, column),
		Raw:  rawType,
		Null: column.Nullable,
	}
	atlasColumn.SetNull(column.Nullable)
	if column.DefaultValue != "" {
		atlasColumn.SetDefault(&atlasschema.RawExpr{X: column.DefaultValue})
	}
	if column.AutoIncrement {
		atlasAddAutoIncrementAttr(dialectName, atlasColumn)
	}
	return atlasColumn
}

func atlasColumnRawType(dialectName string, column ColumnMeta) string {
	rawType := strings.TrimSpace(column.SQLType)
	if rawType == "" {
		rawType = inferTypeName(column)
	}
	if strings.EqualFold(strings.TrimSpace(dialectName), "sqlite") && column.AutoIncrement {
		return "integer"
	}
	return rawType
}

func atlasColumnType(driver atlasmigrate.Driver, rawType string, column ColumnMeta) atlasschema.Type {
	if parser, ok := driver.(atlasschema.TypeParser); ok && rawType != "" {
		if parsed, err := parser.ParseType(rawType); err == nil && parsed != nil {
			return parsed
		}
	}
	return atlasFallbackType(rawType, column)
}

func (c *atlasCompiledSchema) atlasForeignKeyForSpec(table *atlasschema.Table, foreignKey ForeignKeyMeta) *atlasschema.ForeignKey {
	columns := lo.FilterMap(foreignKey.Columns, func(name string, _ int) (*atlasschema.Column, bool) {
		column, ok := table.Column(name)
		return column, ok
	})
	if len(columns) == 0 {
		return nil
	}
	refTable := c.referenceTable(table.Schema, foreignKey.TargetTable, foreignKey.TargetColumns)
	refColumns := lo.FilterMap(foreignKey.TargetColumns, func(name string, _ int) (*atlasschema.Column, bool) {
		column, ok := refTable.Column(name)
		return column, ok
	})
	if len(refColumns) == 0 {
		return nil
	}
	return atlasschema.NewForeignKey(foreignKey.Name).
		SetTable(table).
		AddColumns(columns...).
		SetRefTable(refTable).
		AddRefColumns(refColumns...).
		SetOnDelete(atlasReferenceAction(foreignKey.OnDelete)).
		SetOnUpdate(atlasReferenceAction(foreignKey.OnUpdate))
}

func (c *atlasCompiledSchema) referenceTable(schema *atlasschema.Schema, tableName string, targetColumns []string) *atlasschema.Table {
	if compiled, ok := c.tables.Get(tableName); ok {
		return compiled.table
	}
	if external, ok := c.externals.Get(tableName); ok {
		return external
	}
	external := atlasschema.NewTable(tableName).SetSchema(schema)
	for i := range targetColumns {
		external.AddColumns(atlasschema.NewColumn(targetColumns[i]))
	}
	c.externals.Set(tableName, external)
	return external
}

func atlasReferenceAction(action ReferentialAction) atlasschema.ReferenceOption {
	switch normalized := normalizeReferentialAction(action); normalized {
	case ReferentialCascade:
		return atlasschema.Cascade
	case ReferentialRestrict:
		return atlasschema.Restrict
	case ReferentialSetNull:
		return atlasschema.SetNull
	case ReferentialSetDefault:
		return atlasschema.SetDefault
	case ReferentialNoAction:
		return atlasschema.NoAction
	default:
		return atlasschema.NoAction
	}
}
