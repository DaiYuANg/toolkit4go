// Package main demonstrates basic dbx CRUD and transaction flows.
package main

import (
	"context"
	"fmt"

	"github.com/DaiYuANg/arcgo/dbx"
	"github.com/DaiYuANg/arcgo/examples/dbx/internal/shared"
)

func main() {
	ctx := context.Background()
	catalog := shared.NewCatalog()

	core, closeDB := openBasicDB()
	defer closeOrPanic(closeDB)

	prepareBasicData(ctx, core, catalog)

	printActiveUsers(queryActiveUsers(ctx, core, catalog))
	printUserSummaries(queryUserSummaries(ctx, core, catalog))
	updateUserStatus(ctx, core, catalog, "bob", 2)
	printUpdatedStatus(queryUsersByUsername(ctx, core, catalog, "bob"))
	printLine("basic example completed")
}

func openBasicDB() (*dbx.DB, func() error) {
	core, closeDB, err := shared.OpenSQLite(
		"dbx-basic",
		dbx.WithLogger(shared.NewLogger()),
		dbx.WithDebug(true),
		dbx.WithHooks(dbx.HookFuncs{
			AfterFunc: func(_ context.Context, event *dbx.HookEvent) {
				if event.Operation == dbx.OperationAutoMigrate && event.Err == nil {
					printLine("hook: auto_migrate finished")
				}
			},
		}),
	)
	if err != nil {
		panic(err)
	}

	return core, closeDB
}

func prepareBasicData(ctx context.Context, core *dbx.DB, catalog shared.Catalog) {
	_, err := core.AutoMigrate(ctx, catalog.Roles, catalog.Users, catalog.UserRoles)
	if err != nil {
		panic(err)
	}
	err = shared.SeedDemoData(ctx, core, catalog)
	if err != nil {
		panic(err)
	}
}

func queryActiveUsers(ctx context.Context, core *dbx.DB, catalog shared.Catalog) []shared.User {
	userMapper := dbx.MustMapper[shared.User](catalog.Users)
	users, err := dbx.QueryAll[shared.User](
		ctx,
		core,
		dbx.Select(catalog.Users.AllColumns().Values()...).
			From(catalog.Users).
			Where(catalog.Users.Status.Eq(1)).
			OrderBy(catalog.Users.ID.Asc()),
		userMapper,
	)
	if err != nil {
		panic(err)
	}

	return users
}

func printActiveUsers(users []shared.User) {
	printLine("active users:")
	for index := range users {
		user := &users[index]
		printFormat("- id=%d username=%s email=%s role_id=%d\n", user.ID, user.Username, user.Email, user.RoleID)
	}
}

func queryUserSummaries(ctx context.Context, core *dbx.DB, catalog shared.Catalog) []shared.UserSummary {
	summaryMapper := dbx.MustMapper[shared.UserSummary](catalog.Users)
	summaries, err := dbx.QueryAll[shared.UserSummary](
		ctx,
		core,
		dbx.MustSelectMapped(catalog.Users, summaryMapper).OrderBy(catalog.Users.ID.Asc()),
		summaryMapper,
	)
	if err != nil {
		panic(err)
	}

	return summaries
}

func printUserSummaries(summaries []shared.UserSummary) {
	printLine("projected summaries:")
	for index := range summaries {
		summary := &summaries[index]
		printFormat("- id=%d username=%s email=%s\n", summary.ID, summary.Username, summary.Email)
	}
}

func updateUserStatus(ctx context.Context, core *dbx.DB, catalog shared.Catalog, username string, status int) {
	tx, err := core.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}

	_, err = dbx.Exec(
		ctx,
		tx,
		dbx.Update(catalog.Users).
			Set(catalog.Users.Status.Set(status)).
			Where(catalog.Users.Username.Eq(username)),
	)
	if err != nil {
		rollbackOrPanic(tx.Rollback)
		panic(err)
	}

	//nolint:contextcheck // dbx.Tx commit API does not accept context.
	commitOrPanic(tx)
}

func queryUsersByUsername(ctx context.Context, core *dbx.DB, catalog shared.Catalog, username string) []shared.User {
	userMapper := dbx.MustMapper[shared.User](catalog.Users)
	users, err := dbx.QueryAll[shared.User](
		ctx,
		core,
		dbx.Select(catalog.Users.AllColumns().Values()...).
			From(catalog.Users).
			Where(catalog.Users.Username.Eq(username)),
		userMapper,
	)
	if err != nil {
		panic(err)
	}

	return users
}

func printUpdatedStatus(users []shared.User) {
	printFormat("bob status after tx update: %d\n", users[0].Status)
}

func rollbackOrPanic(rollback func() error) {
	err := rollback()
	if err != nil {
		panic(err)
	}
}

func commitOrPanic(tx *dbx.Tx) {
	err := tx.Commit()
	if err != nil {
		panic(err)
	}
}

func closeOrPanic(closeFn func() error) {
	err := closeFn()
	if err != nil {
		panic(err)
	}
}

func printLine(text string) {
	if _, err := fmt.Println(text); err != nil {
		panic(err)
	}
}

func printFormat(format string, args ...any) {
	if _, err := fmt.Printf(format, args...); err != nil {
		panic(err)
	}
}
