package dbx

import (
	"context"

	"github.com/segmentio/ksuid"
)

type ksuidGenerator struct{}

func NewKSUIDGenerator() IDGenerator {
	return ksuidGenerator{}
}

func (ksuidGenerator) GenerateID(_ context.Context, column ColumnMeta) (any, error) {
	if column.IDStrategy != IDStrategyKSUID {
		return nil, unsupportedIDStrategy(column.IDStrategy)
	}
	return ksuid.New().String(), nil
}
