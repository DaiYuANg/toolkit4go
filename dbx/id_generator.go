package dbx

import (
	"context"
	"fmt"
	"hash/fnv"
	"os"
)

type IDGenerator interface {
	GenerateID(ctx context.Context, column ColumnMeta) (any, error)
}

const (
	DefaultNodeID uint16 = 1
	MinNodeID     uint16 = 1
	MaxNodeID     uint16 = 1023
)

func ResolveNodeIDFromHostName() uint16 {
	hostName, err := os.Hostname()
	if err != nil || hostName == "" {
		return DefaultNodeID
	}
	hasher := fnv.New32a()
	if _, err := hasher.Write([]byte(hostName)); err != nil {
		return DefaultNodeID
	}
	id := uint16(hasher.Sum32() % (uint32(MaxNodeID) + 1))
	if id < MinNodeID {
		return MinNodeID
	}
	return id
}

func unsupportedIDStrategy(strategy IDStrategy) error {
	return fmt.Errorf("dbx: unsupported id strategy %q", strategy)
}
