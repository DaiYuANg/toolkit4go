package dbx

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type uuidGenerator struct{}

func NewUUIDGenerator() IDGenerator {
	return uuidGenerator{}
}

func (uuidGenerator) GenerateID(_ context.Context, column ColumnMeta) (any, error) {
	if column.IDStrategy != IDStrategyUUID {
		return nil, unsupportedIDStrategy(column.IDStrategy)
	}
	return nextUUID(column.UUIDVersion)
}

func nextUUID(version string) (string, error) {
	switch version {
	case "", "v7":
		id, err := uuid.NewV7()
		if err != nil {
			return "", wrapDBError("generate uuid v7", err)
		}
		return id.String(), nil
	case "v4":
		return uuid.NewString(), nil
	default:
		return "", fmt.Errorf("dbx: unsupported uuid version %q", version)
	}
}
