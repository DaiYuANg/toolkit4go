package dbx

import (
	"context"
	"testing"
)

func TestDefaultIDGeneratorDispatchesStrategies(t *testing.T) {
	generator, err := NewDefaultIDGenerator(DefaultNodeID)
	if err != nil {
		t.Fatalf("NewDefaultIDGenerator returned error: %v", err)
	}

	cases := []struct {
		name     string
		column   ColumnMeta
		validate func(any) bool
	}{
		{
			name:     "snowflake",
			column:   ColumnMeta{IDStrategy: IDStrategySnowflake},
			validate: func(value any) bool { id, ok := value.(int64); return ok && id > 0 },
		},
		{
			name:     "uuid",
			column:   ColumnMeta{IDStrategy: IDStrategyUUID, UUIDVersion: DefaultUUIDVersion},
			validate: func(value any) bool { id, ok := value.(string); return ok && id != "" },
		},
		{
			name:     "ulid",
			column:   ColumnMeta{IDStrategy: IDStrategyULID},
			validate: func(value any) bool { id, ok := value.(string); return ok && id != "" },
		},
		{
			name:     "ksuid",
			column:   ColumnMeta{IDStrategy: IDStrategyKSUID},
			validate: func(value any) bool { id, ok := value.(string); return ok && id != "" },
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			id, err := generator.GenerateID(context.Background(), tt.column)
			if err != nil {
				t.Fatalf("GenerateID returned error: %v", err)
			}
			if !tt.validate(id) {
				t.Fatalf("unexpected generated id: %#v", id)
			}
		})
	}
}

func TestSnowflakeGeneratorRejectsOtherStrategies(t *testing.T) {
	generator, err := NewSnowflakeGenerator(DefaultNodeID)
	if err != nil {
		t.Fatalf("NewSnowflakeGenerator returned error: %v", err)
	}

	if _, err := generator.GenerateID(context.Background(), ColumnMeta{IDStrategy: IDStrategyUUID}); err == nil {
		t.Fatal("expected snowflake generator to reject uuid strategy")
	}
}
