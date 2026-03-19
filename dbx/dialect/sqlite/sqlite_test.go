package sqlite

import (
	"context"
	"database/sql/driver"
	"reflect"
	"testing"

	"github.com/DaiYuANg/arcgo/dbx"
	"github.com/DaiYuANg/arcgo/dbx/internal/testsql"
)

func TestBuildCreateTable(t *testing.T) {
	bound, err := Dialect{}.BuildCreateTable(dbx.TableSpec{
		Name:        "users",
		Columns:     []dbx.ColumnMeta{{Name: "id", Table: "users", GoType: reflect.TypeFor[int64](), PrimaryKey: true, AutoIncrement: true}, {Name: "username", Table: "users", GoType: reflect.TypeFor[string]()}, {Name: "email_address", Table: "users", GoType: reflect.TypeFor[string]()}, {Name: "role_id", Table: "users", GoType: reflect.TypeFor[int64]()}, {Name: "status", Table: "users", GoType: reflect.TypeFor[int]()}},
		PrimaryKey:  &dbx.PrimaryKeyMeta{Name: "pk_users", Table: "users", Columns: []string{"id"}},
		ForeignKeys: []dbx.ForeignKeyMeta{{Name: "fk_users_role_id", Table: "users", Columns: []string{"role_id"}, TargetTable: "roles", TargetColumns: []string{"id"}, OnDelete: dbx.ReferentialCascade}},
		Checks:      []dbx.CheckMeta{{Name: "ck_users_status", Table: "users", Expression: "status >= 0"}},
	})
	if err != nil {
		t.Fatalf("BuildCreateTable returned error: %v", err)
	}
	expected := `CREATE TABLE IF NOT EXISTS "users" ("id" INTEGER PRIMARY KEY AUTOINCREMENT, "username" TEXT NOT NULL, "email_address" TEXT NOT NULL, "role_id" INTEGER NOT NULL, "status" INTEGER NOT NULL, CONSTRAINT "fk_users_role_id" FOREIGN KEY ("role_id") REFERENCES "roles" ("id") ON DELETE CASCADE, CONSTRAINT "ck_users_status" CHECK (status >= 0))`
	if bound.SQL != expected {
		t.Fatalf("unexpected create table sql:\nwant: %s\n got: %s", expected, bound.SQL)
	}
}

func TestInspectTable(t *testing.T) {
	sqlDB, _, cleanup, err := testsql.Open(testsql.Plan{
		Queries: []testsql.QueryPlan{
			{SQL: "SELECT name FROM sqlite_master WHERE type = 'table' AND name = ?", Args: []driver.Value{"users"}, Columns: []string{"name"}, Rows: [][]driver.Value{{"users"}}},
			{SQL: `PRAGMA table_info("users")`, Columns: []string{"cid", "name", "type", "notnull", "dflt_value", "pk"}, Rows: [][]driver.Value{{int64(0), "id", "INTEGER", int64(1), nil, int64(1)}, {int64(1), "username", "TEXT", int64(1), nil, int64(0)}, {int64(2), "email_address", "TEXT", int64(1), nil, int64(0)}, {int64(3), "role_id", "INTEGER", int64(1), nil, int64(0)}, {int64(4), "status", "INTEGER", int64(1), nil, int64(0)}}},
			{SQL: `PRAGMA index_list("users")`, Columns: []string{"seq", "name", "unique", "origin", "partial"}, Rows: [][]driver.Value{{int64(0), "idx_users_username", int64(0), "c", int64(0)}, {int64(1), "ux_users_email_address", int64(1), "c", int64(0)}}},
			{SQL: `PRAGMA index_info("idx_users_username")`, Columns: []string{"seqno", "cid", "name"}, Rows: [][]driver.Value{{int64(0), int64(1), "username"}}},
			{SQL: `PRAGMA index_info("ux_users_email_address")`, Columns: []string{"seqno", "cid", "name"}, Rows: [][]driver.Value{{int64(0), int64(2), "email_address"}}},
			{SQL: `PRAGMA foreign_key_list("users")`, Columns: []string{"id", "seq", "table", "from", "to", "on_update", "on_delete", "match"}, Rows: [][]driver.Value{{int64(0), int64(0), "roles", "role_id", "id", "NO ACTION", "CASCADE", "NONE"}}},
			{SQL: "SELECT sql FROM sqlite_master WHERE type = 'table' AND name = ?", Args: []driver.Value{"users"}, Columns: []string{"sql"}, Rows: [][]driver.Value{{`CREATE TABLE users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT NOT NULL, email_address TEXT NOT NULL, role_id INTEGER NOT NULL, status INTEGER NOT NULL, CONSTRAINT ck_users_status CHECK (status >= 0))`}}},
		},
	})
	if err != nil {
		t.Fatalf("testsql.Open returned error: %v", err)
	}
	defer cleanup()

	state, err := Dialect{}.InspectTable(context.Background(), sqlDB, "users")
	if err != nil {
		t.Fatalf("InspectTable returned error: %v", err)
	}
	if !state.Exists || len(state.Columns) != 5 || len(state.Indexes) != 2 {
		t.Fatalf("unexpected table state: %+v", state)
	}
	if state.PrimaryKey == nil || len(state.PrimaryKey.Columns) != 1 || state.PrimaryKey.Columns[0] != "id" {
		t.Fatalf("unexpected primary key state: %+v", state.PrimaryKey)
	}
	if len(state.ForeignKeys) != 1 || state.ForeignKeys[0].TargetTable != "roles" {
		t.Fatalf("unexpected foreign key state: %+v", state.ForeignKeys)
	}
	if len(state.Checks) != 1 || state.Checks[0].Expression != "status >= 0" {
		t.Fatalf("unexpected check state: %+v", state.Checks)
	}
}
