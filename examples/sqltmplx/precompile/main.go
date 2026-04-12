// Package main demonstrates precompiled sqltmplx templates.
package main

import (
	"fmt"

	"github.com/DaiYuANg/arcgo/dbx/dialect/mysql"
	"github.com/DaiYuANg/arcgo/dbx/sqltmplx"
)

func main() {
	engine := sqltmplx.New(mysql.New())
	tpl, err := engine.Compile(`
SELECT id, order_no, status
FROM orders
/*%where */
/*%if present(Tenant) */
  AND tenant = /* Tenant */'acme'
/*%end */
/*%if present(Status) */
  AND status = /* Status */'PAID'
/*%end */
/*%end */
ORDER BY id DESC
LIMIT /* Page.Limit */20 OFFSET /* Page.Offset */0
`)
	if err != nil {
		panic(err)
	}

	for _, params := range []map[string]any{
		{"Tenant": "acme", "Status": "PAID"},
		{"Tenant": "acme", "Status": "SHIPPED"},
	} {
		bound, err := tpl.RenderPage(params, sqltmplx.Page(1, 20))
		if err != nil {
			panic(err)
		}

		mustPrintln(bound.Query)
		mustPrintln(bound.Args)
	}
}

func mustPrintln(value any) {
	if _, err := fmt.Println(value); err != nil {
		panic(err)
	}
}
