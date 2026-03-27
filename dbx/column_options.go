package dbx

import "strings"

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

func WithReference[E any, T any](ref ForeignKeyRef) ColumnOption[E, T] {
	return func(column Column[E, T]) Column[E, T] {
		column.meta.References = new(ref)
		return column
	}
}

func WithIDStrategyColumn[E any, T any](strategy IDStrategy) ColumnOption[E, T] {
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
	return WithIDStrategyColumn[E, T](IDStrategyDBAuto)
}

func SnowflakeIDColumn[E any, T any]() ColumnOption[E, T] {
	return WithIDStrategyColumn[E, T](IDStrategySnowflake)
}

func UUIDIDColumn[E any, T any]() ColumnOption[E, T] {
	return WithIDStrategyColumn[E, T](IDStrategyUUID)
}

func UUIDv7IDColumn[E any, T any]() ColumnOption[E, T] {
	return func(column Column[E, T]) Column[E, T] {
		column.meta.IDStrategy = IDStrategyUUID
		column.meta.UUIDVersion = "v7"
		return column
	}
}

func UUIDv4IDColumn[E any, T any]() ColumnOption[E, T] {
	return func(column Column[E, T]) Column[E, T] {
		column.meta.IDStrategy = IDStrategyUUID
		column.meta.UUIDVersion = "v4"
		return column
	}
}

func (c Column[E, T]) bindColumn(binding columnBinding) any {
	meta := c.meta
	meta.Name = binding.meta.Name
	meta.Table = binding.meta.Table
	meta.Alias = binding.meta.Alias
	meta.FieldName = binding.meta.FieldName
	if meta.GoType == nil {
		meta.GoType = binding.meta.GoType
	}
	if meta.SQLType == "" {
		meta.SQLType = binding.meta.SQLType
	}
	meta.PrimaryKey = meta.PrimaryKey || binding.meta.PrimaryKey
	if meta.IDStrategy == IDStrategyUnset {
		meta.AutoIncrement = meta.AutoIncrement || binding.meta.AutoIncrement
	} else {
		meta.AutoIncrement = meta.IDStrategy == IDStrategyDBAuto
	}
	meta.Nullable = meta.Nullable || binding.meta.Nullable
	meta.Unique = meta.Unique || binding.meta.Unique
	meta.Indexed = meta.Indexed || binding.meta.Indexed
	if meta.DefaultValue == "" {
		meta.DefaultValue = binding.meta.DefaultValue
	}
	if meta.References == nil && binding.meta.References != nil {
		meta.References = new(*binding.meta.References)
	}
	if meta.IDStrategy == IDStrategyUnset {
		meta.IDStrategy = binding.meta.IDStrategy
	}
	if meta.UUIDVersion == "" {
		meta.UUIDVersion = binding.meta.UUIDVersion
	}
	if meta.IDStrategy == IDStrategyUUID && meta.UUIDVersion == "" {
		meta.UUIDVersion = DefaultUUIDVersion
	}
	c.meta = meta
	return c
}
