package dbx

import (
	schemax "github.com/DaiYuANg/arcgo/dbx/schema"
	"strings"

	"github.com/DaiYuANg/arcgo/dbx/idgen"
)

func PrimaryKeyColumn[E any, T any]() ColumnOption[E, T] {
	return func(column Column[E, T]) Column[E, T] {
		column.meta.PrimaryKey = true
		return column
	}
}

func AutoIncrementColumn[E any, T any]() ColumnOption[E, T] {
	return func(column Column[E, T]) Column[E, T] {
		column.meta.AutoIncrement = true
		return column
	}
}

func NullableColumn[E any, T any]() ColumnOption[E, T] {
	return func(column Column[E, T]) Column[E, T] {
		column.meta.Nullable = true
		return column
	}
}

func UniqueColumn[E any, T any]() ColumnOption[E, T] {
	return func(column Column[E, T]) Column[E, T] {
		column.meta.Unique = true
		return column
	}
}

func IndexedColumn[E any, T any]() ColumnOption[E, T] {
	return func(column Column[E, T]) Column[E, T] {
		column.meta.Indexed = true
		return column
	}
}

func WithDefault[E any, T any](value string) ColumnOption[E, T] {
	return func(column Column[E, T]) Column[E, T] {
		column.meta.DefaultValue = value
		return column
	}
}

func WithReference[E any, T any](ref schemax.ForeignKeyRef) ColumnOption[E, T] {
	return func(column Column[E, T]) Column[E, T] {
		column.meta.References = new(ref)
		return column
	}
}

func WithIDStrategyColumn[E any, T any](strategy idgen.Strategy) ColumnOption[E, T] {
	return func(column Column[E, T]) Column[E, T] {
		column.meta.IDStrategy = strategy
		return column
	}
}

func WithUUIDVersionColumn[E any, T any](version string) ColumnOption[E, T] {
	return func(column Column[E, T]) Column[E, T] {
		column.meta.UUIDVersion = strings.TrimSpace(version)
		return column
	}
}

func DBAutoIDColumn[E any, T any]() ColumnOption[E, T] {
	return WithIDStrategyColumn[E, T](idgen.StrategyDBAuto)
}

func SnowflakeIDColumn[E any, T any]() ColumnOption[E, T] {
	return WithIDStrategyColumn[E, T](idgen.StrategySnowflake)
}

func UUIDIDColumn[E any, T any]() ColumnOption[E, T] {
	return WithIDStrategyColumn[E, T](idgen.StrategyUUID)
}

func UUIDv7IDColumn[E any, T any]() ColumnOption[E, T] {
	return func(column Column[E, T]) Column[E, T] {
		column.meta.IDStrategy = idgen.StrategyUUID
		column.meta.UUIDVersion = "v7"
		return column
	}
}

func UUIDv4IDColumn[E any, T any]() ColumnOption[E, T] {
	return func(column Column[E, T]) Column[E, T] {
		column.meta.IDStrategy = idgen.StrategyUUID
		column.meta.UUIDVersion = "v4"
		return column
	}
}

func (c Column[E, T]) bindColumn(binding columnBinding) any {
	meta := c.meta
	mergeColumnBasic(&meta, binding.meta)
	mergeColumnFlags(&meta, binding.meta)
	mergeColumnDefaultsAndRefs(&meta, binding.meta)
	finalizeColumnIDAndUUID(&meta, binding.meta)
	c.meta = meta
	return c
}

func mergeColumnBasic(meta *schemax.ColumnMeta, b schemax.ColumnMeta) {
	meta.Name = b.Name
	meta.Table = b.Table
	meta.Alias = b.Alias
	meta.FieldName = b.FieldName
	if meta.GoType == nil {
		meta.GoType = b.GoType
	}
	if meta.SQLType == "" {
		meta.SQLType = b.SQLType
	}
}

func mergeColumnFlags(meta *schemax.ColumnMeta, b schemax.ColumnMeta) {
	meta.PrimaryKey = meta.PrimaryKey || b.PrimaryKey
	if meta.IDStrategy == idgen.StrategyUnset {
		meta.AutoIncrement = meta.AutoIncrement || b.AutoIncrement
	} else {
		meta.AutoIncrement = meta.IDStrategy == idgen.StrategyDBAuto
	}
	meta.Nullable = meta.Nullable || b.Nullable
	meta.Unique = meta.Unique || b.Unique
	meta.Indexed = meta.Indexed || b.Indexed
}

func mergeColumnDefaultsAndRefs(meta *schemax.ColumnMeta, b schemax.ColumnMeta) {
	if meta.DefaultValue == "" {
		meta.DefaultValue = b.DefaultValue
	}
	if meta.References == nil && b.References != nil {
		meta.References = new(*b.References)
	}
}

func finalizeColumnIDAndUUID(meta *schemax.ColumnMeta, b schemax.ColumnMeta) {
	if meta.IDStrategy == idgen.StrategyUnset {
		meta.IDStrategy = b.IDStrategy
	}
	if meta.UUIDVersion == "" {
		meta.UUIDVersion = b.UUIDVersion
	}
	if meta.IDStrategy == idgen.StrategyUUID && meta.UUIDVersion == "" {
		meta.UUIDVersion = idgen.DefaultUUIDVersion
	}
}
