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

	core, closeDB, err := shared.OpenSQLite("dbx-migration", dbx.WithLogger(shared.NewLogger()), dbx.WithDebug(true))
	if err != nil {
		panic(err)
	}
	defer func() { _ = closeDB() }()

	plan, err := core.PlanSchemaChanges(ctx, catalog.Roles, catalog.Users, catalog.UserRoles)
	if err != nil {
		panic(err)
	}

	fmt.Println("planned migration actions:")
	for _, action := range plan.Actions {
		fmt.Printf("- kind=%s executable=%t summary=%s\n", action.Kind, action.Executable, action.Summary)
		if action.Statement.SQL != "" {
			fmt.Printf("  sql=%s\n", action.Statement.SQL)
		}
	}

	report, err := core.AutoMigrate(ctx, catalog.Roles, catalog.Users, catalog.UserRoles)
	if err != nil {
		panic(err)
	}
	fmt.Printf("auto migrate valid=%t tables=%d\n", report.Valid(), len(report.Tables))

	validated, err := core.ValidateSchemas(ctx, catalog.Roles, catalog.Users, catalog.UserRoles)
	if err != nil {
		panic(err)
	}
	fmt.Printf("validate valid=%t\n", validated.Valid())

	fmt.Println("users foreign keys:")
	for _, fk := range catalog.Users.ForeignKeys() {
		fmt.Printf("- name=%s columns=%v target=%s(%v)\n", fk.Name, fk.Columns, fk.TargetTable, fk.TargetColumns)
	}
}
