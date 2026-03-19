package repository

import (
	"database/sql"
	"testing"

	"github.com/DaiYuANg/arcgo/dbx"
	sqlitedialect "github.com/DaiYuANg/arcgo/dbx/dialect/sqlite"
)

type User struct {
	ID   int64  `dbx:"id"`
	Name string `dbx:"name"`
}

type UserSchema struct {
	dbx.Schema[User]
	ID   dbx.Column[User, int64]  `dbx:"id,pk,auto"`
	Name dbx.Column[User, string] `dbx:"name"`
}

func TestNewUsesSchemaAsMetadataSource(t *testing.T) {
	core := dbx.New((*sql.DB)(nil), sqlitedialect.Dialect{})
	users := dbx.MustSchema("users", UserSchema{})
	repo := New[User](core, users)

	if repo.DB() != core {
		t.Fatal("expected repository to hold db core")
	}
	if repo.Schema().TableName() != "users" {
		t.Fatalf("unexpected schema table: %q", repo.Schema().TableName())
	}
	if _, ok := repo.Mapper().FieldByColumn("name"); !ok {
		t.Fatal("expected mapper to expose name column")
	}
}
