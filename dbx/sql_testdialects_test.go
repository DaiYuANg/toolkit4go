package dbx

import (
	"fmt"
	"strconv"
	"strings"
)

type testPostgresDialect struct{}

func (testPostgresDialect) Name() string { return "postgres" }
func (testPostgresDialect) BindVar(n int) string {
	return "$" + strconv.Itoa(n)
}
func (testPostgresDialect) QuoteIdent(ident string) string {
	return `"` + strings.ReplaceAll(ident, `"`, `""`) + `"`
}
func (testPostgresDialect) RenderLimitOffset(limit, offset *int) (string, error) {
	if limit == nil && offset == nil {
		return "", nil
	}
	if limit != nil && offset != nil {
		return fmt.Sprintf("LIMIT %d OFFSET %d", *limit, *offset), nil
	}
	if limit != nil {
		return fmt.Sprintf("LIMIT %d", *limit), nil
	}
	return fmt.Sprintf("OFFSET %d", *offset), nil
}

type testMySQLDialect struct{}

func (testMySQLDialect) Name() string         { return "mysql" }
func (testMySQLDialect) BindVar(_ int) string { return "?" }
func (testMySQLDialect) QuoteIdent(ident string) string {
	return "`" + strings.ReplaceAll(ident, "`", "``") + "`"
}
func (testMySQLDialect) RenderLimitOffset(limit, offset *int) (string, error) {
	if limit == nil && offset == nil {
		return "", nil
	}
	if limit != nil && offset != nil {
		return fmt.Sprintf("LIMIT %d OFFSET %d", *limit, *offset), nil
	}
	if limit != nil {
		return fmt.Sprintf("LIMIT %d", *limit), nil
	}
	return fmt.Sprintf("LIMIT 18446744073709551615 OFFSET %d", *offset), nil
}
