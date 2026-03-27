// Package main demonstrates sqltmplx update rendering with SQLite syntax.
package main

import (
	"fmt"

	"github.com/DaiYuANg/arcgo/dbx/dialect/sqlite"
	"github.com/DaiYuANg/arcgo/dbx/sqltmplx"
)

// UpdateCommand contains the fields used by the SQLite update template.
type UpdateCommand struct {
	ID     int
	Name   string `db:"name"`
	Status string `json:"status"`
}

func main() {
	engine := sqltmplx.New(sqlite.New())

	tpl := `
UPDATE users
/*%set */
/*%if present(Name) */
  name = /* Name */'alice',
/*%end */
/*%if present(Status) */
  status = /* Status */'active',
/*%end */
/*%end */
WHERE id = /* ID */1
`

	bound, err := engine.Render(tpl, UpdateCommand{
		ID:     42,
		Name:   "alice",
		Status: "active",
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
