// Package main demonstrates sqltmplx rendering with PostgreSQL syntax.
package main

import (
	"fmt"

	"github.com/DaiYuANg/arcgo/dbx/dialect/postgres"
	"github.com/DaiYuANg/arcgo/dbx/sqltmplx"
	"github.com/DaiYuANg/arcgo/dbx/sqltmplx/validate"
	_ "github.com/DaiYuANg/arcgo/dbx/sqltmplx/validate/postgresparser"
)

// Query contains the parameters rendered into the PostgreSQL template.
type Query struct {
	Tenant string `db:"tenant"`
	Name   string `json:"name"`
	IDs    []int  `json:"ids"`
}

func main() {
	engine := sqltmplx.New(
		postgres.New(),
		sqltmplx.WithValidator(validate.NewSQLParser(postgres.New())),
	)

	tpl := `
SELECT id, tenant, name
FROM users
/*%where */
/*%if present(Tenant) */
  AND tenant = /* Tenant */'acme'
/*%end */
/*%if present(Name) */
  AND name = /* Name */'alice'
/*%end */
/*%if !empty(IDs) */
  AND id IN (/* IDs */(1, 2, 3))
/*%end */
/*%end */
ORDER BY id DESC
`

	bound, err := engine.Render(tpl, Query{
		Tenant: "acme",
		Name:   "alice",
		IDs:    []int{1, 2, 3},
	})
	if err != nil {
		panic(err)
	}

	mustPrintln(bound.Query)
	mustPrintln(bound.Args)
}

func mustPrintln(value any) {
	if _, err := fmt.Println(value); err != nil {
		panic(err)
	}
}
