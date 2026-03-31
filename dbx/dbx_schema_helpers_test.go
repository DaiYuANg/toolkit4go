package dbx_test

import (
	"testing"
)

func assertUserSchemaBasics(t *testing.T, users UserSchema) {
	t.Helper()
	if users.TableName() != "users" {
		t.Fatalf("unexpected table name: %q", users.TableName())
	}
	if users.ID.Ref() != "users.id" {
		t.Fatalf("unexpected id ref: %q", users.ID.Ref())
	}
	if users.Email.Name() != "email_address" {
		t.Fatalf("unexpected email column name: %q", users.Email.Name())
	}
	if !users.ID.IsPrimaryKey() || !users.ID.Meta().AutoIncrement {
		t.Fatalf("expected id metadata to mark pk/auto: %+v", users.ID.Meta())
	}
	if users.ID.Meta().IDStrategy != IDStrategyDBAuto {
		t.Fatalf("expected default int64 id strategy db_auto, got: %q", users.ID.Meta().IDStrategy)
	}
}

func assertUserReferenceMetadata(t *testing.T, users UserSchema) {
	t.Helper()
	ref, ok := users.RoleID.Reference()
	if !ok {
		t.Fatal("expected role_id reference metadata")
	}
	if ref.TargetTable != "roles" || ref.TargetColumn != "id" || ref.OnDelete != ReferentialCascade {
		t.Fatalf("unexpected reference metadata: %+v", ref)
	}
}

func assertUserRelationMetadata(t *testing.T, users UserSchema) {
	t.Helper()
	columns := users.Columns()
	if len(columns) != 5 {
		t.Fatalf("unexpected columns metadata count: %d", len(columns))
	}
	relations := users.Relations()
	if len(relations) != 4 {
		t.Fatalf("unexpected relations metadata count: %d", len(relations))
	}
	if relations[0].Kind != RelationBelongsTo || relations[0].TargetTable != "roles" {
		t.Fatalf("unexpected first relation metadata: %+v", relations[0])
	}
	if relations[3].Kind != RelationManyToMany || relations[3].ThroughTable != "user_roles" {
		t.Fatalf("unexpected many-to-many metadata: %+v", relations[3])
	}
}

func assertUserForeignKeys(t *testing.T, users UserSchema) {
	t.Helper()
	foreignKeys := users.ForeignKeys()
	if len(foreignKeys) != 1 {
		t.Fatalf("unexpected foreign key count: %d", len(foreignKeys))
	}
	if foreignKeys[0].Columns[0] != "role_id" || foreignKeys[0].TargetTable != "roles" {
		t.Fatalf("unexpected foreign key metadata: %+v", foreignKeys[0])
	}
}
