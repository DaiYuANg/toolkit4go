package mysql

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
	expected := "CREATE TABLE IF NOT EXISTS `users` (`id` BIGINT AUTO_INCREMENT PRIMARY KEY, `username` TEXT NOT NULL, `email_address` TEXT NOT NULL, `role_id` BIGINT NOT NULL, `status` INT NOT NULL, CONSTRAINT `fk_users_role_id` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE, CONSTRAINT `ck_users_status` CHECK (status >= 0))"
	if bound.SQL != expected {
		t.Fatalf("unexpected create sql:\nwant: %s\n got: %s", expected, bound.SQL)
	}
}

func TestInspectTable(t *testing.T) {
	sqlDB, _, cleanup, err := testsql.Open(testsql.Plan{Queries: []testsql.QueryPlan{
		{SQL: "SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?", Args: []driver.Value{"users"}, Columns: []string{"table_name"}, Rows: [][]driver.Value{{"users"}}},
		{SQL: "SELECT column_name, column_type, is_nullable, column_default, column_key, extra FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = ? ORDER BY ordinal_position", Args: []driver.Value{"users"}, Columns: []string{"column_name", "column_type", "is_nullable", "column_default", "column_key", "extra"}, Rows: [][]driver.Value{{"id", "bigint", "NO", nil, "PRI", "auto_increment"}, {"username", "text", "NO", nil, "", ""}, {"email_address", "text", "NO", nil, "", ""}, {"role_id", "bigint", "NO", nil, "", ""}, {"status", "int", "NO", nil, "", ""}}},
		{SQL: "SELECT index_name, non_unique, column_name FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = ? ORDER BY index_name, seq_in_index", Args: []driver.Value{"users"}, Columns: []string{"index_name", "non_unique", "column_name"}, Rows: [][]driver.Value{{"idx_users_username", int64(1), "username"}, {"ux_users_email_address", int64(0), "email_address"}}},
		{SQL: "SELECT kcu.constraint_name, kcu.column_name, kcu.referenced_table_name, kcu.referenced_column_name, rc.UPDATE_RULE, rc.DELETE_RULE FROM information_schema.key_column_usage kcu JOIN information_schema.table_constraints tc ON kcu.constraint_name = tc.constraint_name AND kcu.table_schema = tc.table_schema AND kcu.table_name = tc.table_name LEFT JOIN information_schema.referential_constraints rc ON kcu.constraint_name = rc.constraint_name AND kcu.table_schema = rc.constraint_schema WHERE kcu.table_schema = DATABASE() AND kcu.table_name = ? AND tc.constraint_type = 'FOREIGN KEY' ORDER BY kcu.constraint_name, kcu.ordinal_position", Args: []driver.Value{"users"}, Columns: []string{"constraint_name", "column_name", "referenced_table_name", "referenced_column_name", "UPDATE_RULE", "DELETE_RULE"}, Rows: [][]driver.Value{{"fk_users_role_id", "role_id", "roles", "id", "NO ACTION", "CASCADE"}}},
		{SQL: "SELECT tc.constraint_name, cc.check_clause FROM information_schema.table_constraints tc JOIN information_schema.check_constraints cc ON tc.constraint_name = cc.constraint_name AND tc.constraint_schema = cc.constraint_schema WHERE tc.table_schema = DATABASE() AND tc.table_name = ? AND tc.constraint_type = 'CHECK' ORDER BY tc.constraint_name", Args: []driver.Value{"users"}, Columns: []string{"constraint_name", "check_clause"}, Rows: [][]driver.Value{{"ck_users_status", "status >= 0"}}},
	}})
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
	if state.PrimaryKey == nil || state.PrimaryKey.Name != "PRIMARY" {
		t.Fatalf("unexpected primary key state: %+v", state.PrimaryKey)
	}
	if len(state.ForeignKeys) != 1 || state.ForeignKeys[0].TargetTable != "roles" {
		t.Fatalf("unexpected foreign key state: %+v", state.ForeignKeys)
	}
	if len(state.Checks) != 1 || state.Checks[0].Expression != "status >= 0" {
		t.Fatalf("unexpected check state: %+v", state.Checks)
	}
}
