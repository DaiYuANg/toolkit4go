package dbx

import (
	"context"
	"database/sql"
	"strings"
	"testing"
)

func BenchmarkCompileAtlasSchema(b *testing.B) {
	roles := MustSchema("roles", advancedRoleSchema{})
	users := MustSchema("users", advancedUserSchema{})
	schemas := []SchemaResource{roles, users}

	b.ReportAllocs()
	for range b.N {
		if compiled := compileAtlasSchema("sqlite", nil, "main", schemas); compiled == nil {
			b.Fatal("compileAtlasSchema returned nil")
		}
	}
}

func BenchmarkPlanSchemaChangesSQLiteAtlasEmpty(b *testing.B) {
	ctx := context.Background()
	roles := MustSchema("roles", RoleSchema{})
	users := MustSchema("users", UserSchema{})

	run := func(b *testing.B, db *sql.DB) {
		b.Helper()
		core := New(db, testSQLiteDialect{})
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			if _, err := core.PlanSchemaChanges(ctx, roles, users); err != nil {
				b.Fatalf("PlanSchemaChanges returned error: %v", err)
			}
		}
	}

	b.Run("Memory", func(b *testing.B) {
		run(b, benchmarkOpenSQLiteDBMemory(b))
	})
	b.Run("IO", func(b *testing.B) {
		run(b, benchmarkOpenSQLiteDB(b, "plan-empty"))
	})
}

func BenchmarkValidateSchemasSQLiteAtlasMatched(b *testing.B) {
	ctx := context.Background()
	roles := MustSchema("roles", RoleSchema{})
	users := MustSchema("users", UserSchema{})

	run := func(b *testing.B, db *sql.DB) {
		b.Helper()
		core := New(db, testSQLiteDialect{})
		if _, err := core.AutoMigrate(ctx, roles, users); err != nil {
			b.Fatalf("AutoMigrate returned error: %v", err)
		}
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			if _, err := core.ValidateSchemas(ctx, roles, users); err != nil {
				b.Fatalf("ValidateSchemas returned error: %v", err)
			}
		}
	}

	b.Run("Memory", func(b *testing.B) {
		run(b, benchmarkOpenSQLiteDBMemory(b))
	})
	b.Run("IO", func(b *testing.B) {
		run(b, benchmarkOpenSQLiteDB(b, "validate-matched"))
	})
}

func BenchmarkMigrationPlanSQLPreview(b *testing.B) {
	ctx := context.Background()
	roles := MustSchema("roles", RoleSchema{})
	users := MustSchema("users", UserSchema{})

	run := func(b *testing.B, db *sql.DB) {
		b.Helper()
		core := New(db, testSQLiteDialect{})
		plan, err := core.PlanSchemaChanges(ctx, roles, users)
		if err != nil {
			b.Fatalf("PlanSchemaChanges returned error: %v", err)
		}
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			preview := plan.SQLPreview()
			if len(preview) == 0 || !strings.Contains(strings.ToLower(preview[0]), "create table") {
				b.Fatalf("unexpected preview: %+v", preview)
			}
		}
	}

	b.Run("Memory", func(b *testing.B) {
		run(b, benchmarkOpenSQLiteDBMemory(b))
	})
	b.Run("IO", func(b *testing.B) {
		run(b, benchmarkOpenSQLiteDB(b, "preview"))
	})
}
