package dbx

import (
	"fmt"

	"github.com/DaiYuANg/arcgo/collectionx"
	"github.com/samber/lo"
)

func ProjectionOf[E any](schema SchemaResource, mapper Mapper[E]) ([]SelectItem, error) {
	return projectionOfDefinition(schema.schemaRef(), mapper)
}

func MustProjectionOf[E any](schema SchemaResource, mapper Mapper[E]) []SelectItem {
	items, err := projectionOfDefinition(schema.schemaRef(), mapper)
	if err != nil {
		panic(err)
	}
	return items
}

func SelectMapped[E any](schema SchemaResource, mapper Mapper[E]) (*SelectQuery, error) {
	items, err := projectionOfDefinition(schema.schemaRef(), mapper)
	if err != nil {
		return nil, err
	}
	return Select(items...).From(schema), nil
}

func MustSelectMapped[E any](schema SchemaResource, mapper Mapper[E]) *SelectQuery {
	items, err := projectionOfDefinition(schema.schemaRef(), mapper)
	if err != nil {
		panic(err)
	}
	return Select(items...).From(schema)
}

func projectionOfDefinition[E any](definition schemaDefinition, mapper Mapper[E]) ([]SelectItem, error) {
	columns := lo.Associate(definition.columns, func(column ColumnMeta) (string, ColumnMeta) {
		return column.Name, column
	})

	items := collectionx.NewListWithCapacity[SelectItem](mapper.meta.fields.Len())
	mapper.meta.fields.Range(func(_ int, field MappedField) bool {
		column, ok := columns[field.Column]
		if !ok {
			return false
		}
		items.Add(schemaSelectItem{meta: column})
		return true
	})
	if items.Len() != mapper.meta.fields.Len() {
		for _, field := range mapper.meta.fields.Values() {
			if _, ok := columns[field.Column]; !ok {
				return nil, fmt.Errorf("%w: %s", ErrUnmappedColumn, field.Column)
			}
		}
	}
	return items.Values(), nil
}
