package idgen

import (
	"context"

	"github.com/DaiYuANg/arcgo/collectionx"
)

type defaultGenerator struct {
	generators collectionx.Map[Strategy, Generator]
}

func NewDefault(nodeID uint16) (Generator, error) {
	snowflake, err := NewSnowflake(nodeID)
	if err != nil {
		return nil, err
	}

	generators := collectionx.NewMapWithCapacity[Strategy, Generator](4)
	generators.Set(StrategySnowflake, snowflake)
	generators.Set(StrategyUUID, NewUUID())
	generators.Set(StrategyULID, NewULID())
	generators.Set(StrategyKSUID, NewKSUID())

	return &defaultGenerator{generators: generators}, nil
}

func (g *defaultGenerator) GenerateID(ctx context.Context, request Request) (any, error) {
	if g == nil {
		return nil, unsupportedStrategy(request.Strategy)
	}
	generator, ok := g.generators.Get(request.Strategy)
	if !ok {
		return nil, unsupportedStrategy(request.Strategy)
	}
	return generator.GenerateID(ctx, request)
}
