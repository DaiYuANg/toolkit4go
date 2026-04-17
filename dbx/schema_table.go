package dbx

import (
	"reflect"
	"strings"

	"github.com/DaiYuANg/arcgo/collectionx"
	"github.com/DaiYuANg/arcgo/dbx/querydsl"
)

type SchemaSource[E any] interface {
	querydsl.TableSource
	schemaRef() schemaDefinition
}

type schemaBinder interface {
	bindSchema(def schemaDefinition) any
	entityType() reflect.Type
}

type schemaDefinition struct {
	table         querydsl.Table
	columns       collectionx.List[ColumnMeta]
	columnsByName collectionx.Map[string, ColumnMeta]
	relations     collectionx.List[RelationMeta]
	indexes       collectionx.List[IndexMeta]
	primaryKey    *PrimaryKeyMeta
	checks        collectionx.List[CheckMeta]
}

type Schema[E any] struct {
	def schemaDefinition
}

func (s Schema[E]) bindSchema(def schemaDefinition) any {
	s.def = def
	return s
}

func (Schema[E]) entityType() reflect.Type {
	return reflect.TypeFor[E]()
}

func (s Schema[E]) schemaRef() schemaDefinition {
	return s.def
}

func (s Schema[E]) tableRef() querydsl.Table {
	return s.def.table
}

func (s Schema[E]) Name() string {
	return s.def.table.Name()
}

func (s Schema[E]) TableName() string {
	return s.def.table.Name()
}

func (s Schema[E]) Alias() string {
	return s.def.table.Alias()
}

func (s Schema[E]) TableAlias() string {
	return s.def.table.Alias()
}

func (s Schema[E]) Ref() string {
	return s.def.table.Ref()
}

func (s Schema[E]) QualifiedName() string {
	return s.def.table.QualifiedName()
}

func (s Schema[E]) EntityType() reflect.Type {
	return s.def.table.EntityType()
}

func (s Schema[E]) Columns() collectionx.List[ColumnMeta] {
	return cloneColumnMetas(s.def.columns)
}

func (s Schema[E]) Relations() collectionx.List[RelationMeta] {
	return s.def.relations.Clone()
}

func (s Schema[E]) Indexes() collectionx.List[IndexMeta] {
	return cloneIndexMetas(s.def.indexes)
}

func (s Schema[E]) PrimaryKey() (PrimaryKeyMeta, bool) {
	if s.def.primaryKey == nil {
		return PrimaryKeyMeta{}, false
	}
	return clonePrimaryKeyMeta(*s.def.primaryKey), true
}

func (s Schema[E]) Checks() collectionx.List[CheckMeta] {
	return cloneCheckMetas(s.def.checks)
}

func (s Schema[E]) ForeignKeys() collectionx.List[ForeignKeyMeta] {
	items := deriveForeignKeys(s.def)
	return collectionx.MapList(collectionx.NewListWithCapacity(len(items), items...), func(_ int, item ForeignKeyMeta) ForeignKeyMeta {
		return cloneForeignKeyMeta(item)
	})
}

func MustSchema[S any](name string, schema S) S {
	bound, err := bindSchema(name, "", schema)
	if err != nil {
		panic(err)
	}
	return bound
}

func Alias[S querydsl.TableSource](schema S, alias string) S {
	if strings.TrimSpace(alias) == "" {
		panic("dbx: alias cannot be empty")
	}
	bound, err := bindSchema(tableRef(schema).Name(), alias, schema)
	if err != nil {
		panic(err)
	}
	return bound
}

func tableRef(source querydsl.TableSource) querydsl.Table {
	if source == nil {
		return querydsl.Table{}
	}
	if typed, ok := source.(interface{ tableRef() querydsl.Table }); ok {
		return typed.tableRef()
	}
	return querydsl.NewTableRef(source.TableName(), source.TableAlias(), nil, nil)
}

func (d schemaDefinition) columnByName(name string) (ColumnMeta, bool) {
	if d.columnsByName != nil && d.columnsByName.Len() > 0 {
		return d.columnsByName.Get(name)
	}
	return collectionx.FindList(d.columns, func(_ int, column ColumnMeta) bool {
		return column.Name == name
	})
}

func cloneColumnMetas(items collectionx.List[ColumnMeta]) collectionx.List[ColumnMeta] {
	return collectionx.MapList(items, func(_ int, column ColumnMeta) ColumnMeta {
		return cloneColumnMeta(column)
	})
}

func cloneIndexMetas(items collectionx.List[IndexMeta]) collectionx.List[IndexMeta] {
	return collectionx.MapList(items, func(_ int, item IndexMeta) IndexMeta {
		return cloneIndexMeta(item)
	})
}

func cloneCheckMetas(items collectionx.List[CheckMeta]) collectionx.List[CheckMeta] {
	return collectionx.MapList(items, func(_ int, item CheckMeta) CheckMeta {
		return cloneCheckMeta(item)
	})
}

func indexColumnsByName(columns collectionx.List[ColumnMeta]) collectionx.Map[string, ColumnMeta] {
	return collectionx.AssociateList(columns, func(_ int, column ColumnMeta) (string, ColumnMeta) {
		return column.Name, column
	})
}
