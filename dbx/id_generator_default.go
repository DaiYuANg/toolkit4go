package dbx

import (
	"context"

	"github.com/DaiYuANg/arcgo/collectionx"
)

type defaultIDGenerator struct {
	generators collectionx.Map[IDStrategy, IDGenerator]
}

func NewDefaultIDGenerator(nodeID uint16) (IDGenerator, error) {
	snowflake, err := NewSnowflakeGenerator(nodeID)
	if err != nil {
		return nil, err
	}

	generators := collectionx.NewMapWithCapacity[IDStrategy, IDGenerator](4)
	generators.Set(IDStrategySnowflake, snowflake)
	generators.Set(IDStrategyUUID, NewUUIDGenerator())
	generators.Set(IDStrategyULID, NewULIDGenerator())
	generators.Set(IDStrategyKSUID, NewKSUIDGenerator())

	return &defaultIDGenerator{generators: generators}, nil
}

func (g *defaultIDGenerator) GenerateID(ctx context.Context, column ColumnMeta) (any, error) {
	if g == nil {
		return nil, unsupportedIDStrategy(column.IDStrategy)
	}
	generator, ok := g.generators.Get(column.IDStrategy)
	if !ok {
		return nil, unsupportedIDStrategy(column.IDStrategy)
	}
	return generator.GenerateID(ctx, column)
}
