package dbx

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/DaiYuANg/arcgo/collectionx"
	"github.com/DaiYuANg/arcgo/dbx/querydsl"
	relationx "github.com/DaiYuANg/arcgo/dbx/relation"
	schemax "github.com/DaiYuANg/arcgo/dbx/schema"
)

type schemaBindingState struct {
	binder        schemaBinder
	binderField   reflect.Value
	defTable      querydsl.Table
	columns       collectionx.List[schemax.ColumnMeta]
	columnsByName collectionx.Map[string, schemax.ColumnMeta]
	relations     collectionx.List[schemax.RelationMeta]
	indexes       collectionx.List[schemax.IndexMeta]
	checks        collectionx.List[schemax.CheckMeta]
	primaryKey    *schemax.PrimaryKeyMeta
}

func bindSchema[S any](name, alias string, schema S) (S, error) {
	value := reflect.ValueOf(&schema).Elem()
	if value.Kind() != reflect.Struct {
		return schema, errors.New("dbx: schema must be a struct")
	}

	schemaType := value.Type()
	state := newSchemaBindingState(schemaType, name, alias, value.NumField())
	for i := range value.NumField() {
		if err := state.bindField(schemaType.Field(i), value.Field(i)); err != nil {
			return schema, err
		}
	}

	def, err := state.definition(schemaType)
	if err != nil {
		return schema, err
	}
	state.binderField.Set(reflect.ValueOf(state.binder.bindSchema(def)))
	return schema, nil
}

func newSchemaBindingState(schemaType reflect.Type, name, alias string, fieldCount int) schemaBindingState {
	return schemaBindingState{
		defTable:      querydsl.NewTableRef(strings.TrimSpace(name), strings.TrimSpace(alias), schemaType, nil),
		columns:       collectionx.NewListWithCapacity[schemax.ColumnMeta](fieldCount),
		columnsByName: collectionx.NewMapWithCapacity[string, schemax.ColumnMeta](fieldCount),
		relations:     collectionx.NewListWithCapacity[schemax.RelationMeta](fieldCount),
		indexes:       collectionx.NewListWithCapacity[schemax.IndexMeta](fieldCount),
		checks:        collectionx.NewListWithCapacity[schemax.CheckMeta](fieldCount),
	}
}

func (s *schemaBindingState) bindField(fieldType reflect.StructField, fieldValue reflect.Value) error {
	if !fieldValue.CanSet() {
		return nil
	}
	if s.captureSchemaBinder(fieldValue) {
		return nil
	}
	if handled, err := s.bindColumnField(fieldType, fieldValue); handled || err != nil {
		return err
	}
	if s.bindRelationField(fieldType, fieldValue) {
		return nil
	}
	if handled, err := s.bindConstraintField(fieldType, fieldValue); handled || err != nil {
		return err
	}
	return nil
}

func (s *schemaBindingState) captureSchemaBinder(fieldValue reflect.Value) bool {
	candidate, ok := fieldValue.Interface().(schemaBinder)
	if !ok {
		return false
	}
	s.binder = candidate
	s.binderField = fieldValue
	s.defTable = s.defTable.WithEntityType(candidate.entityType())
	return true
}

func (s *schemaBindingState) bindColumnField(fieldType reflect.StructField, fieldValue reflect.Value) (bool, error) {
	candidate, ok := fieldValue.Interface().(columnBinder)
	if !ok {
		return false, nil
	}
	meta, err := resolveColumnMeta(s.defTable, fieldType, fieldValue.Interface())
	if err != nil {
		return true, err
	}
	bound := candidate.bindColumn(columnBinding{meta: meta})
	fieldValue.Set(reflect.ValueOf(bound))
	column := meta
	if accessor, ok := bound.(columnAccessor); ok {
		column = accessor.columnRef()
	}
	column = cloneColumnMeta(column)
	s.columns.Add(column)
	s.columnsByName.Set(column.Name, column)
	return true, nil
}

func (s *schemaBindingState) bindRelationField(fieldType reflect.StructField, fieldValue reflect.Value) bool {
	candidate, ok := fieldValue.Interface().(relationx.Binder)
	if !ok {
		return false
	}
	meta := resolveRelationMeta(s.defTable, fieldType, candidate)
	fieldValue.Set(reflect.ValueOf(candidate.BindRelation(relationx.Binding{Meta: meta})))
	s.relations.Add(meta)
	return true
}

func (s *schemaBindingState) bindConstraintField(fieldType reflect.StructField, fieldValue reflect.Value) (bool, error) {
	candidate, ok := fieldValue.Interface().(constraintBinder)
	if !ok {
		return false, nil
	}
	binding, err := resolveConstraintBinding(s.defTable, fieldType, fieldValue.Interface())
	if err != nil {
		return true, err
	}
	fieldValue.Set(reflect.ValueOf(candidate.bindConstraint(binding)))
	s.indexes.Add(binding.indexes...)
	if binding.primaryKey != nil {
		s.primaryKey = new(clonePrimaryKeyMeta(*binding.primaryKey))
	}
	if binding.check != nil {
		s.checks.Add(*binding.check)
	}
	return true, nil
}

func (s *schemaBindingState) definition(schemaType reflect.Type) (schemaDefinition, error) {
	if s.binder == nil {
		return schemaDefinition{}, fmt.Errorf("dbx: schema %s must embed dbx.Schema[T]", schemaType.Name())
	}
	columns := cloneColumnMetas(s.columns)
	columnsByName := s.columnsByName.Clone()
	if columnsByName.Len() == 0 && columns.Len() > 0 {
		columnsByName = indexColumnsByName(columns)
	}
	return schemaDefinition{
		table:         s.defTable,
		columns:       columns,
		columnsByName: columnsByName,
		relations:     s.relations.Clone(),
		indexes:       cloneIndexMetas(s.indexes),
		primaryKey:    clonePrimaryKeyMetaPtr(s.primaryKey),
		checks:        cloneCheckMetas(s.checks),
	}, nil
}

func resolveColumnMeta(def querydsl.Table, field reflect.StructField, value any) (schemax.ColumnMeta, error) {
	name, options := resolveTagNameAndOptions(field)
	meta := schemax.ColumnMeta{
		Name:          name,
		Table:         def.Name(),
		Alias:         def.Alias(),
		FieldName:     field.Name,
		GoType:        resolveColumnGoType(value),
		SQLType:       optionValue(options, "type"),
		PrimaryKey:    optionEnabled(options, "pk"),
		AutoIncrement: optionEnabled(options, "auto") || optionEnabled(options, "autoincrement"),
		Nullable:      optionEnabled(options, "nullable") || optionEnabled(options, "null"),
		Unique:        optionEnabled(options, "unique"),
		Indexed:       optionEnabled(options, "index") || optionEnabled(options, "indexed"),
		DefaultValue:  optionValue(options, "default"),
	}

	if refValue := optionValue(options, "ref"); refValue != "" {
		targetTable, targetColumn, ok := splitReference(refValue)
		if ok {
			meta.References = &schemax.ForeignKeyRef{
				TargetTable:  targetTable,
				TargetColumn: targetColumn,
				OnDelete:     parseReferentialAction(optionValue(options, "ondelete")),
				OnUpdate:     parseReferentialAction(optionValue(options, "onupdate")),
			}
		}
	}

	return normalizeIDPolicy(meta)
}

func resolveColumnGoType(value any) reflect.Type {
	reporter, ok := value.(columnTypeReporter)
	if !ok {
		return nil
	}
	return reporter.valueType()
}

func resolveRelationMeta(def querydsl.Table, field reflect.StructField, binder relationx.Binder) schemax.RelationMeta {
	options := parseTagOptions(field.Tag.Get("rel"))
	name := optionValue(options, "name")
	if name == "" {
		name = toSnakeCase(field.Name)
	}
	return schemax.RelationMeta{
		Name:                name,
		FieldName:           field.Name,
		Kind:                binder.RelationKind(),
		SourceTable:         def.Name(),
		SourceAlias:         def.Alias(),
		TargetTable:         optionValue(options, "table"),
		LocalColumn:         optionValue(options, "local"),
		TargetColumn:        optionValue(options, "target"),
		ThroughTable:        optionValue(options, "join"),
		ThroughLocalColumn:  optionValue(options, "join_local"),
		ThroughTargetColumn: optionValue(options, "join_target"),
		TargetType:          binder.TargetType(),
	}
}
