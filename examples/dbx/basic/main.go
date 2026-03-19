package main

import (
	"context"
	"fmt"
	"os"

	"github.com/DaiYuANg/arcgo/dbx"
	"github.com/DaiYuANg/arcgo/examples/dbx/internal/shared"
)

func main() {
	ctx := context.Background()
	catalog := shared.NewCatalog()
	logger := shared.NewLogger()

	core, closeDB, err := shared.OpenSQLite(
		"dbx-basic",
		dbx.WithLogger(logger),
		dbx.WithDebug(true),
		dbx.WithHooks(dbx.HookFuncs{
			AfterFunc: func(_ context.Context, event *dbx.HookEvent) {
				if event.Operation == dbx.OperationAutoMigrate && event.Err == nil {
					fmt.Println("hook: auto_migrate finished")
				}
			},
		}),
	)
	if err != nil {
		panic(err)
	}
	defer func() { _ = closeDB() }()

	if _, err := core.AutoMigrate(ctx, catalog.Roles, catalog.Users, catalog.UserRoles); err != nil {
		panic(err)
	}
	if err := shared.SeedDemoData(ctx, core, catalog); err != nil {
		panic(err)
	}

	userMapper := dbx.MustMapper[shared.User](catalog.Users)
	activeUsers, err := dbx.QueryAll(
		ctx,
		core,
		dbx.Select(catalog.Users.AllColumns()...).
			From(catalog.Users).
			Where(catalog.Users.Status.Eq(1)).
			OrderBy(catalog.Users.ID.Asc()),
		userMapper,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("active users:")
	for _, user := range activeUsers {
		fmt.Printf("- id=%d username=%s email=%s role_id=%d\n", user.ID, user.Username, user.Email, user.RoleID)
	}

	summaryMapper := dbx.MustMapper[shared.UserSummary](catalog.Users)
	summaries, err := dbx.QueryAll(
		ctx,
		core,
		dbx.MustSelectMapped(catalog.Users, summaryMapper).OrderBy(catalog.Users.ID.Asc()),
		summaryMapper,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("projected summaries:")
	for _, summary := range summaries {
		fmt.Printf("- id=%d username=%s email=%s\n", summary.ID, summary.Username, summary.Email)
	}

	tx, err := core.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}
	if _, err := dbx.Exec(ctx, tx, dbx.Update(catalog.Users).Set(catalog.Users.Status.Set(2)).Where(catalog.Users.Username.Eq("bob"))); err != nil {
		_ = tx.Rollback()
		panic(err)
	}
	if err := tx.Commit(); err != nil {
		panic(err)
	}

	updated, err := dbx.QueryAll(
		ctx,
		core,
		dbx.Select(catalog.Users.AllColumns()...).
			From(catalog.Users).
			Where(catalog.Users.Username.Eq("bob")),
		userMapper,
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("bob status after tx update: %d\n", updated[0].Status)

	_, _ = fmt.Fprintln(os.Stdout, "basic example completed")
}
