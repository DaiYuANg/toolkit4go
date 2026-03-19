package dbx

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/DaiYuANg/arcgo/collectionx"
)

type Mapper[E any] struct {
	meta *mapperMetadata
}

type MappedField struct {
	Name       string
	Column     string
	Index      int
	Insertable bool
	Updatable  bool
}

type scanPlan struct {
	fields []MappedField
}

type mapperMetadata struct {
	entityType reflect.Type
	fields     collectionx.List[MappedField]
	byColumn   collectionx.Map[string, MappedField]
	scanPlans  collectionx.ConcurrentMap[string, *scanPlan]
}

type mapperCacheKey struct {
	entityType reflect.Type
	schemaType reflect.Type
}

var mapperCache = collectionx.NewConcurrentMap[mapperCacheKey, *mapperMetadata]()

func MustMapper[E any, S SchemaSource[E]](schema S) Mapper[E] {
	meta, err := getOrBuildMapperMetadata[E](schema)
	if err != nil {
		panic(err)
	}
	return Mapper[E]{meta: meta}
}

func (m Mapper[E]) Fields() []MappedField {
	return m.meta.fields.Values()
}

func (m Mapper[E]) FieldByColumn(column string) (MappedField, bool) {
	return m.meta.byColumn.Get(column)
}

func (m Mapper[E]) ScanRows(rows *sql.Rows) ([]E, error) {
	if rows == nil {
		return nil, fmt.Errorf("dbx: rows is nil")
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	plan, err := m.scanPlan(columns)
	if err != nil {
		return nil, err
	}

	items := collectionx.NewList[E]()
	for rows.Next() {
		entity, err := m.scanCurrentRow(rows, plan)
		if err != nil {
			return nil, err
		}
		items.Add(entity)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items.Values(), nil
}

func (m Mapper[E]) InsertAssignments(schema SchemaSource[E], entity *E) ([]Assignment, error) {
	return m.entityAssignments(schema, entity, func(column ColumnMeta, field MappedField) bool {
		if !field.Insertable {
			return false
		}
		return !(column.PrimaryKey && column.AutoIncrement)
	})
}

func (m Mapper[E]) UpdateAssignments(schema SchemaSource[E], entity *E) ([]Assignment, error) {
	return m.entityAssignments(schema, entity, func(column ColumnMeta, field MappedField) bool {
		if !field.Updatable {
			return false
		}
		return !column.PrimaryKey && !column.AutoIncrement
	})
}

func (m Mapper[E]) PrimaryPredicate(schema SchemaSource[E], entity *E) (Predicate, error) {
	value, err := m.entityValue(entity)
	if err != nil {
		return nil, err
	}

	for _, column := range schema.schemaRef().columns {
		if !column.PrimaryKey {
			continue
		}
		field, ok := m.meta.byColumn.Get(column.Name)
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrPrimaryKeyUnmapped, column.Name)
		}
		return metadataComparisonPredicate{
			left:  column,
			op:    OpEq,
			right: value.Field(field.Index).Interface(),
		}, nil
	}

	return nil, ErrNoPrimaryKey
}

func getOrBuildMapperMetadata[E any, S SchemaSource[E]](schema S) (*mapperMetadata, error) {
	meta := schema.schemaRef()
	entityType := reflect.TypeFor[E]()
	key := mapperCacheKey{entityType: entityType, schemaType: meta.table.schemaType}
	if cached, ok := mapperCache.Get(key); ok {
		return cached, nil
	}

	mapper, err := buildMapperMetadata(entityType, meta.columns)
	if err != nil {
		return nil, err
	}
	actual, _ := mapperCache.GetOrStore(key, mapper)
	return actual, nil
}

func buildMapperMetadata(entityType reflect.Type, columns []ColumnMeta) (*mapperMetadata, error) {
	if entityType.Kind() != reflect.Struct {
		return nil, ErrUnsupportedEntity
	}

	columnSet := collectionx.NewSetWithCapacity[string](len(columns))
	for _, column := range columns {
		columnSet.Add(column.Name)
	}

	fields := collectionx.NewListWithCapacity[MappedField](entityType.NumField())
	byColumn := collectionx.NewMapWithCapacity[string, MappedField](len(columns))
	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		if !field.IsExported() {
			continue
		}

		columnName, options := resolveEntityColumn(field)
		if columnName == "" || !columnSet.Contains(columnName) {
			continue
		}

		mapped := MappedField{
			Name:       field.Name,
			Column:     columnName,
			Index:      i,
			Insertable: !options["readonly"] && !options["-insert"] && !options["noinsert"],
			Updatable:  !options["readonly"] && !options["-update"] && !options["noupdate"],
		}
		fields.Add(mapped)
		byColumn.Set(columnName, mapped)
	}

	return &mapperMetadata{
		entityType: entityType,
		fields:     fields,
		byColumn:   byColumn,
		scanPlans:  collectionx.NewConcurrentMapWithCapacity[string, *scanPlan](8),
	}, nil
}

func resolveEntityColumn(field reflect.StructField) (string, map[string]bool) {
	raw := strings.TrimSpace(field.Tag.Get("dbx"))
	if raw == "-" {
		return "", nil
	}
	if raw == "" {
		return toSnakeCase(field.Name), map[string]bool{}
	}

	parts := strings.Split(raw, ",")
	name := strings.TrimSpace(parts[0])
	if name == "" {
		name = toSnakeCase(field.Name)
	}
	options := make(map[string]bool, len(parts)-1)
	for _, option := range parts[1:] {
		trimmed := strings.ToLower(strings.TrimSpace(option))
		if trimmed == "" {
			continue
		}
		options[trimmed] = true
	}
	return name, options
}

func (m Mapper[E]) scanPlan(columns []string) (*scanPlan, error) {
	signature := scanSignature(columns)
	if cached, ok := m.meta.scanPlans.Get(signature); ok {
		return cached, nil
	}

	fields := collectionx.NewListWithCapacity[MappedField](len(columns))
	for _, column := range columns {
		field, ok := m.meta.byColumn.Get(column)
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnmappedColumn, column)
		}
		fields.Add(field)
	}

	plan := &scanPlan{fields: fields.Values()}
	actual, _ := m.meta.scanPlans.GetOrStore(signature, plan)
	return actual, nil
}

func (m Mapper[E]) scanCurrentRow(rows *sql.Rows, plan *scanPlan) (E, error) {
	value := reflect.New(m.meta.entityType).Elem()
	destinations := make([]any, len(plan.fields))
	for i, field := range plan.fields {
		destinations[i] = value.Field(field.Index).Addr().Interface()
	}

	if err := rows.Scan(destinations...); err != nil {
		var zero E
		return zero, err
	}
	return value.Interface().(E), nil
}

func (m Mapper[E]) entityAssignments(schema SchemaSource[E], entity *E, include func(column ColumnMeta, field MappedField) bool) ([]Assignment, error) {
	value, err := m.entityValue(entity)
	if err != nil {
		return nil, err
	}

	assignments := collectionx.NewListWithCapacity[Assignment](len(schema.schemaRef().columns))
	for _, column := range schema.schemaRef().columns {
		field, ok := m.meta.byColumn.Get(column.Name)
		if !ok || !include(column, field) {
			continue
		}
		assignments.Add(metadataAssignment{
			meta:  column,
			value: value.Field(field.Index).Interface(),
		})
	}

	return assignments.Values(), nil
}

func (m Mapper[E]) entityValue(entity *E) (reflect.Value, error) {
	if entity == nil {
		return reflect.Value{}, ErrNilEntity
	}
	return reflect.ValueOf(entity).Elem(), nil
}

func scanSignature(columns []string) string {
	return strings.Join(columns, "\x1f")
}
