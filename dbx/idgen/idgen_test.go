package idgen

import (
	"context"
	"testing"
)

func TestDefaultGeneratorDispatchesStrategies(t *testing.T) {
	generator, err := NewDefault(DefaultNodeID)
	if err != nil {
		t.Fatalf("NewDefault returned error: %v", err)
	}

	cases := []struct {
		name     string
		request  Request
		validate func(any) bool
	}{
		{
			name:     "snowflake",
			request:  Request{Strategy: StrategySnowflake},
			validate: func(value any) bool { id, ok := value.(int64); return ok && id > 0 },
		},
		{
			name:     "uuid",
			request:  Request{Strategy: StrategyUUID, UUIDVersion: DefaultUUIDVersion},
			validate: func(value any) bool { id, ok := value.(string); return ok && id != "" },
		},
		{
			name:     "ulid",
			request:  Request{Strategy: StrategyULID},
			validate: func(value any) bool { id, ok := value.(string); return ok && id != "" },
		},
		{
			name:     "ksuid",
			request:  Request{Strategy: StrategyKSUID},
			validate: func(value any) bool { id, ok := value.(string); return ok && id != "" },
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			id, err := generator.GenerateID(context.Background(), tt.request)
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
	generator, err := NewSnowflake(DefaultNodeID)
	if err != nil {
		t.Fatalf("NewSnowflake returned error: %v", err)
	}

	if _, err := generator.GenerateID(context.Background(), Request{Strategy: StrategyUUID}); err == nil {
		t.Fatal("expected snowflake generator to reject uuid strategy")
	}
}
