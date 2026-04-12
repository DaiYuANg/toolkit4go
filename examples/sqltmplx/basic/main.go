// Package main demonstrates basic sqltmplx rendering with MySQL syntax.
package main

import (
	"fmt"

	"github.com/DaiYuANg/arcgo/dbx/dialect/mysql"
	"github.com/DaiYuANg/arcgo/dbx/sqltmplx"
	"github.com/DaiYuANg/arcgo/dbx/sqltmplx/validate"
	_ "github.com/DaiYuANg/arcgo/dbx/sqltmplx/validate/mysqlparser"
)

func main() {
	engine := sqltmplx.New(
		mysql.New(),
		sqltmplx.WithValidator(validate.NewSQLParser(mysql.New())),
	)

	tpl := `
SELECT id, name, status
FROM users
/*%where */
/*%if name != nil */
  AND name = /* name */'alice'
/*%end */
/*%if ids != nil */
  AND id IN (/* ids */(1, 2, 3))
/*%end */
/*%end */
ORDER BY id DESC
LIMIT /* Page.Limit */20 OFFSET /* Page.Offset */0
`

	bound, err := engine.RenderPage(tpl, map[string]any{
		"name": "alice",
		"ids":  []int{1, 2, 3},
	}, sqltmplx.Page(1, 20))
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
