package activerecord_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DaiYuANg/arcgo/dbx"
	activerecord "github.com/DaiYuANg/arcgo/dbx/activerecord"
	sqlitedialect "github.com/DaiYuANg/arcgo/dbx/dialect/sqlite"
	"github.com/DaiYuANg/arcgo/dbx/repository"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

type User struct {
	ID   int64  `dbx:"id"`
	Name string `dbx:"name"`
}

type UserSchema struct {
	dbx.Schema[User]
	ID   dbx.IDColumn[User, int64, dbx.IDSnowflake] `dbx:"id,pk"`
	Name dbx.Column[User, string]                   `dbx:"name"`
}

func TestModelSaveReloadDelete(t *testing.T) {
	ctx, store := openUserStore(t, "file:activerecord_model_test?mode=memory&cache=shared")

	model := store.Wrap(&User{Name: "alice"})
	require.NoError(t, model.Save(ctx))
	require.NotZero(t, model.Entity().ID)

	model.Entity().Name = "alice-v2"
	require.NoError(t, model.Save(ctx))

	found, err := store.FindByID(ctx, model.Entity().ID)
	require.NoError(t, err)
	require.Equal(t, "alice-v2", found.Entity().Name)

	model.Entity().Name = "stale"
	require.NoError(t, model.Reload(ctx))
	require.Equal(t, "alice-v2", model.Entity().Name)

	require.NoError(t, model.Delete(ctx))

	_, err = store.FindByID(ctx, model.Entity().ID)
	require.True(t, errors.Is(err, repository.ErrNotFound))
}

func TestStoreFindOptionAPIs(t *testing.T) {
	ctx, store := openUserStore(t, "file:activerecord_option_test?mode=memory&cache=shared")

	model := store.Wrap(&User{Name: "alice"})
	require.NoError(t, model.Save(ctx))

	noneByID, err := store.FindByIDOption(ctx, int64(99999))
	require.NoError(t, err)
	require.False(t, noneByID.IsPresent())

	byID, err := store.FindByIDOption(ctx, model.Entity().ID)
	require.NoError(t, err)

	found, ok := byID.Get()
	require.True(t, ok)
	require.Equal(t, "alice", found.Entity().Name)

	byKey, err := store.FindByKeyOption(ctx, found.Key())
	require.NoError(t, err)

	again, ok := byKey.Get()
	require.True(t, ok)
	require.Equal(t, model.Entity().ID, again.Entity().ID)
}

func openUserStore(tb testing.TB, dsn string) (context.Context, *activerecord.Store[User, UserSchema]) {
	tb.Helper()

	ctx := context.Background()
	raw, err := sql.Open("sqlite", dsn)
	require.NoError(tb, err)

	tb.Cleanup(func() {
		if closeErr := raw.Close(); closeErr != nil {
			tb.Errorf("close sqlite: %v", closeErr)
		}
	})

	core := dbx.MustNewWithOptions(raw, sqlitedialect.New())
	users := dbx.MustSchema("users", UserSchema{})

	_, err = core.AutoMigrate(ctx, users)
	require.NoError(tb, err)

	return ctx, activerecord.New[User, UserSchema](core, users)
}
