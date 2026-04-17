package dbx

import (
	"context"

	"github.com/oklog/ulid/v2"
)

type ulidGenerator struct{}

func NewULIDGenerator() IDGenerator {
	return ulidGenerator{}
}

func (ulidGenerator) GenerateID(_ context.Context, column ColumnMeta) (any, error) {
	if column.IDStrategy != IDStrategyULID {
		return nil, unsupportedIDStrategy(column.IDStrategy)
	}
	return ulid.Make().String(), nil
}
