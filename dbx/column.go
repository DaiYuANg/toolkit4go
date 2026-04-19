package dbx

import (
	schemax "github.com/DaiYuANg/arcgo/dbx/schema"
	"reflect"
	"strings"

	"github.com/DaiYuANg/arcgo/dbx/idgen"
	"github.com/DaiYuANg/arcgo/dbx/querydsl"
	"github.com/samber/lo"
)

type Ref[E any, T any] interface {
	Name() string
	refNode()
}

type IDMarker interface {
	idStrategy() idgen.Strategy
	uuidVersion() string
}

func cloneColumnMeta(meta schemax.ColumnMeta) schemax.ColumnMeta {
	if meta.References == nil {
		return meta
	}
	meta.References = new(*meta.References)
	return meta
}

type columnBinder interface {
	bindColumn(binding columnBinding) any
}

type columnAccessor interface {
	columnRef() schemax.ColumnMeta
}

type columnTypeReporter interface {
	valueType() reflect.Type
}

type typedColumn[T any] interface {
	columnRef() schemax.ColumnMeta
}

type columnBinding struct {
	meta schemax.ColumnMeta
}

type Column[E any, T any] struct {
	meta schemax.ColumnMeta
}

// IDColumn declares an ID policy directly in the schema field type.
// The marker strategy is applied during schema binding.
type IDColumn[E any, T any, M IDMarker] struct {
	Column[E, T]
}

type ColumnOption[E any, T any] func(Column[E, T]) Column[E, T]

type IDAuto struct{}
type IDSnowflake struct{}
type IDUUID struct{}
type IDUUIDv7 struct{}
type IDUUIDv4 struct{}
type IDULID struct{}
type IDKSUID struct{}

func (IDAuto) idStrategy() idgen.Strategy { return idgen.StrategyDBAuto }
func (IDAuto) uuidVersion() string        { return "" }
func (IDSnowflake) idStrategy() idgen.Strategy {
	return idgen.StrategySnowflake
}
func (IDSnowflake) uuidVersion() string   { return "" }
func (IDUUID) idStrategy() idgen.Strategy { return idgen.StrategyUUID }
func (IDUUID) uuidVersion() string        { return "" }
func (IDUUIDv7) idStrategy() idgen.Strategy {
	return idgen.StrategyUUID
}
func (IDUUIDv7) uuidVersion() string { return "v7" }
func (IDUUIDv4) idStrategy() idgen.Strategy {
	return idgen.StrategyUUID
}
func (IDUUIDv4) uuidVersion() string       { return "v4" }
func (IDULID) idStrategy() idgen.Strategy  { return idgen.StrategyULID }
func (IDULID) uuidVersion() string         { return "" }
func (IDKSUID) idStrategy() idgen.Strategy { return idgen.StrategyKSUID }
func (IDKSUID) uuidVersion() string        { return "" }

func NewColumn[E any, T any](opts ...ColumnOption[E, T]) Column[E, T] {
	column := Column[E, T]{}
	lo.ForEach(lo.Filter(opts, func(opt ColumnOption[E, T], _ int) bool { return opt != nil }), func(opt ColumnOption[E, T], _ int) {
		column = opt(column)
	})
	return column
}

func (c IDColumn[E, T, M]) bindColumn(binding columnBinding) any {
	marker := *new(M)
	base := c.Column
	base.meta.PrimaryKey = true
	base.meta.IDStrategy = marker.idStrategy()
	if version := marker.uuidVersion(); version != "" {
		base.meta.UUIDVersion = version
	}
	boundValue := base.bindColumn(binding)
	bound, ok := boundValue.(Column[E, T])
	if !ok {
		return IDColumn[E, T, M]{Column: base}
	}
	return IDColumn[E, T, M]{Column: bound}
}

func NamedColumn[T any](source querydsl.TableSource, name string) Column[struct{}, T] {
	table := tableRef(source)
	return Column[struct{}, T]{
		meta: schemax.ColumnMeta{
			Name:   strings.TrimSpace(name),
			Table:  table.Name(),
			Alias:  table.Alias(),
			GoType: reflect.TypeFor[T](),
		},
	}
}

func ResultColumn[T any](name string) Column[struct{}, T] {
	return Column[struct{}, T]{
		meta: schemax.ColumnMeta{
			Name:   strings.TrimSpace(name),
			GoType: reflect.TypeFor[T](),
		},
	}
}
