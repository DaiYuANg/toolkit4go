package sqlite

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/DaiYuANg/arcgo/collectionx"
	"github.com/DaiYuANg/arcgo/dbx"
	"github.com/DaiYuANg/arcgo/dbx/dialect"
	"github.com/samber/lo"
)

var (
	sqliteIntegerKinds = []reflect.Kind{
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
	}
	sqliteRealKinds = []reflect.Kind{reflect.Float32, reflect.Float64}
)

// Dialect implements SQLite rendering and schema inspection.
type Dialect struct{}

// New returns a SQLite dialect implementation.
func New() Dialect { return Dialect{} }

// Name returns the dialect name.
func (Dialect) Name() string { return "sqlite" }

// BindVar returns the bind placeholder for a parameter index.
func (Dialect) BindVar(_ int) string { return "?" }

// QuoteIdent quotes an identifier for SQLite.
func (Dialect) QuoteIdent(ident string) string {
	return `"` + strings.ReplaceAll(ident, `"`, `""`) + `"`
}

// RenderLimitOffset renders a LIMIT/OFFSET clause for SQLite.
func (Dialect) RenderLimitOffset(limit, offset *int) (string, error) {
	if limit == nil && offset == nil {
		return "", nil
	}
	if limit != nil && offset != nil {
		return fmt.Sprintf("LIMIT %d OFFSET %d", *limit, *offset), nil
	}
	if limit != nil {
		return fmt.Sprintf("LIMIT %d", *limit), nil
	}
	return fmt.Sprintf("LIMIT -1 OFFSET %d", *offset), nil
}

// QueryFeatures returns the supported query feature set.
func (Dialect) QueryFeatures() dialect.QueryFeatures {
	return dialect.DefaultQueryFeatures("sqlite")
}

// BuildCreateTable builds a CREATE TABLE statement.
func (d Dialect) BuildCreateTable(spec dbx.TableSpec) (dbx.BoundQuery, error) {
	parts := collectionx.NewListWithCapacity[string](len(spec.Columns) + len(spec.ForeignKeys) + len(spec.Checks) + 1)
	inlinePrimaryKey := singlePrimaryKeyColumn(spec.PrimaryKey)

	for i := range spec.Columns {
		column := spec.Columns[i]
		ddl, err := d.columnDDL(column, columnDDLConfig{
			AllowAutoIncrement: true,
			InlinePrimaryKey:   inlinePrimaryKey == column.Name,
		})
		if err != nil {
			return dbx.BoundQuery{}, fmt.Errorf("build sqlite column ddl: %w", err)
		}
		parts.Add(ddl)
	}

	if spec.PrimaryKey != nil && len(spec.PrimaryKey.Columns) > 1 {
		parts.Add(d.primaryKeyDDL(*spec.PrimaryKey))
	}

	for i := range spec.ForeignKeys {
		parts.Add(d.foreignKeyDDL(spec.ForeignKeys[i]))
	}
	for i := range spec.Checks {
		parts.Add(d.checkDDL(spec.Checks[i]))
	}

	return dbx.BoundQuery{
		SQL: "CREATE TABLE IF NOT EXISTS " + d.QuoteIdent(spec.Name) + " (" + strings.Join(parts.Values(), ", ") + ")",
	}, nil
}

// BuildAddColumn builds an ALTER TABLE ADD COLUMN statement.
func (d Dialect) BuildAddColumn(table string, column dbx.ColumnMeta) (dbx.BoundQuery, error) {
	if column.PrimaryKey {
		return dbx.BoundQuery{}, fmt.Errorf("dbx/sqlite: cannot add primary key column %s with ALTER TABLE", column.Name)
	}

	ddl, err := d.columnDDL(column, columnDDLConfig{IncludeReference: true})
	if err != nil {
		return dbx.BoundQuery{}, fmt.Errorf("build sqlite column ddl: %w", err)
	}

	return dbx.BoundQuery{
		SQL: "ALTER TABLE " + d.QuoteIdent(table) + " ADD COLUMN " + ddl,
	}, nil
}

// BuildCreateIndex builds a CREATE INDEX statement.
func (d Dialect) BuildCreateIndex(index dbx.IndexMeta) (dbx.BoundQuery, error) {
	columns := lo.Map(index.Columns, func(column string, _ int) string {
		return d.QuoteIdent(column)
	})
	prefix := "CREATE INDEX IF NOT EXISTS "
	if index.Unique {
		prefix = "CREATE UNIQUE INDEX IF NOT EXISTS "
	}
	return dbx.BoundQuery{
		SQL: prefix + d.QuoteIdent(index.Name) + " ON " + d.QuoteIdent(index.Table) + " (" + strings.Join(columns, ", ") + ")",
	}, nil
}

// BuildAddForeignKey reports that SQLite foreign keys require a table rebuild.
func (Dialect) BuildAddForeignKey(string, dbx.ForeignKeyMeta) (dbx.BoundQuery, error) {
	return dbx.BoundQuery{}, errors.New("dbx/sqlite: adding foreign keys requires table rebuild")
}

// BuildAddCheck reports that SQLite check constraints require a table rebuild.
func (Dialect) BuildAddCheck(string, dbx.CheckMeta) (dbx.BoundQuery, error) {
	return dbx.BoundQuery{}, errors.New("dbx/sqlite: adding check constraints requires table rebuild")
}

// InspectTable inspects a SQLite table definition from PRAGMA metadata.
func (d Dialect) InspectTable(ctx context.Context, executor dbx.Executor, table string) (dbx.TableState, error) {
	exists, err := inspectSQLiteTableExists(ctx, executor, table)
	if err != nil {
		return dbx.TableState{}, err
	}
	if !exists {
		return dbx.TableState{Name: table, Exists: false}, nil
	}

	columns, primaryKey, err := d.inspectColumns(ctx, executor, table)
	if err != nil {
		return dbx.TableState{}, err
	}

	indexes, err := d.inspectIndexes(ctx, executor, table)
	if err != nil {
		return dbx.TableState{}, err
	}

	foreignKeys, err := d.inspectForeignKeys(ctx, executor, table)
	if err != nil {
		return dbx.TableState{}, err
	}

	checks, autoincrementColumns, err := inspectSQLiteCreateMetadata(ctx, executor, table)
	if err != nil {
		return dbx.TableState{}, err
	}

	return dbx.TableState{
		Exists:      true,
		Name:        table,
		Columns:     markSQLiteAutoincrementColumns(columns, autoincrementColumns),
		Indexes:     indexes,
		PrimaryKey:  primaryKey,
		ForeignKeys: foreignKeys,
		Checks:      checks,
	}, nil
}

// NormalizeType normalizes database type names into dbx logical types.
func (Dialect) NormalizeType(value string) string {
	typeName := strings.ToUpper(strings.TrimSpace(value))
	switch {
	case strings.Contains(typeName, "INT"):
		return "INTEGER"
	case strings.Contains(typeName, "CHAR"), strings.Contains(typeName, "CLOB"), strings.Contains(typeName, "TEXT"):
		return "TEXT"
	case strings.Contains(typeName, "BLOB"):
		return "BLOB"
	case strings.Contains(typeName, "REAL"), strings.Contains(typeName, "FLOA"), strings.Contains(typeName, "DOUBLE"):
		return "REAL"
	case strings.Contains(typeName, "BOOL"):
		return "BOOLEAN"
	case strings.Contains(typeName, "TIMESTAMP"), strings.Contains(typeName, "DATETIME"):
		return "TIMESTAMP"
	default:
		return typeName
	}
}

var _ dbx.SchemaDialect = Dialect{}
