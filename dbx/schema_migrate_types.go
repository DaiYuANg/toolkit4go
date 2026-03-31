package dbx

import (
	"context"
	"strings"

	"github.com/DaiYuANg/arcgo/dbx/dialect"
	"github.com/samber/lo"
)

type SchemaResource interface {
	TableSource
	schemaRef() schemaDefinition
}

type SchemaDialect interface {
	dialect.Dialect
	BuildCreateTable(spec TableSpec) (BoundQuery, error)
	BuildAddColumn(table string, column ColumnMeta) (BoundQuery, error)
	BuildCreateIndex(index IndexMeta) (BoundQuery, error)
	BuildAddForeignKey(table string, foreignKey ForeignKeyMeta) (BoundQuery, error)
	BuildAddCheck(table string, check CheckMeta) (BoundQuery, error)
	InspectTable(ctx context.Context, executor Executor, table string) (TableState, error)
	NormalizeType(value string) string
}

type TableSpec struct {
	Name        string
	Columns     []ColumnMeta
	Indexes     []IndexMeta
	PrimaryKey  *PrimaryKeyMeta
	ForeignKeys []ForeignKeyMeta
	Checks      []CheckMeta
}

type IndexMeta struct {
	Name    string
	Table   string
	Columns []string
	Unique  bool
}

type TableState struct {
	Exists      bool
	Name        string
	Columns     []ColumnState
	Indexes     []IndexState
	PrimaryKey  *PrimaryKeyState
	ForeignKeys []ForeignKeyState
	Checks      []CheckState
}

type ColumnState struct {
	Name          string
	Type          string
	Nullable      bool
	PrimaryKey    bool
	AutoIncrement bool
	DefaultValue  string
}

type IndexState struct {
	Name    string
	Columns []string
	Unique  bool
}

type PrimaryKeyState struct {
	Name    string
	Columns []string
}

type ForeignKeyState struct {
	Name          string
	Columns       []string
	TargetTable   string
	TargetColumns []string
	OnDelete      ReferentialAction
	OnUpdate      ReferentialAction
}

type CheckState struct {
	Name       string
	Expression string
}

type ValidationReport struct {
	Tables   []TableDiff
	Backend  ValidationBackend
	Complete bool
	Warnings []string
}

type ValidationBackend string

const (
	ValidationBackendAtlas  ValidationBackend = "atlas"
	ValidationBackendLegacy ValidationBackend = "legacy"
)

type TableDiff struct {
	Table              string
	MissingTable       bool
	MissingColumns     []ColumnMeta
	MissingIndexes     []IndexMeta
	MissingForeignKeys []ForeignKeyMeta
	MissingChecks      []CheckMeta
	PrimaryKeyDiff     *PrimaryKeyDiff
	ColumnDiffs        []ColumnDiff
}

type PrimaryKeyDiff struct {
	Expected *PrimaryKeyMeta
	Actual   *PrimaryKeyState
	Issues   []string
}

type ColumnDiff struct {
	Column ColumnMeta
	Issues []string
}

type MigrationActionKind string

const (
	MigrationActionCreateTable   MigrationActionKind = "create_table"
	MigrationActionAddColumn     MigrationActionKind = "add_column"
	MigrationActionCreateIndex   MigrationActionKind = "create_index"
	MigrationActionAddForeignKey MigrationActionKind = "add_foreign_key"
	MigrationActionAddCheck      MigrationActionKind = "add_check"
	MigrationActionManual        MigrationActionKind = "manual"
)

type MigrationAction struct {
	Kind       MigrationActionKind
	Table      string
	Summary    string
	Statement  BoundQuery
	Executable bool
}

type MigrationPlan struct {
	Actions []MigrationAction
	Report  ValidationReport
}

func (a MigrationAction) HasStatement() bool {
	return strings.TrimSpace(a.Statement.SQL) != ""
}

func (a MigrationAction) SQLPreview() string {
	return strings.TrimSpace(a.Statement.SQL)
}

func (p MigrationPlan) Statements() []BoundQuery {
	return lo.FilterMap(p.Actions, func(action MigrationAction, _ int) (BoundQuery, bool) {
		return action.Statement, action.HasStatement()
	})
}

func (p MigrationPlan) SQLPreview() []string {
	return lo.FilterMap(p.Actions, func(action MigrationAction, _ int) (string, bool) {
		return action.SQLPreview(), action.HasStatement()
	})
}

type SchemaDriftError struct {
	Report ValidationReport
}

func (e SchemaDriftError) Error() string {
	tables := lo.FilterMap(e.Report.Tables, func(table TableDiff, _ int) (string, bool) {
		return table.Table, !table.Empty()
	})
	if len(tables) == 0 {
		return "dbx: schema drift detected"
	}
	return "dbx: schema drift detected for tables: " + strings.Join(tables, ", ")
}

func (r ValidationReport) Valid() bool {
	return !lo.SomeBy(r.Tables, func(table TableDiff) bool {
		return !table.Empty()
	})
}

func (r ValidationReport) HasWarnings() bool {
	return len(r.Warnings) > 0
}

func (r ValidationReport) IsComplete() bool {
	return r.Complete
}

func (t TableDiff) Empty() bool {
	return !t.MissingTable &&
		len(t.MissingColumns) == 0 &&
		len(t.MissingIndexes) == 0 &&
		len(t.MissingForeignKeys) == 0 &&
		len(t.MissingChecks) == 0 &&
		t.PrimaryKeyDiff == nil &&
		len(t.ColumnDiffs) == 0
}

func (p MigrationPlan) ExecutableActions() []MigrationAction {
	return lo.Filter(p.Actions, func(action MigrationAction, _ int) bool {
		return action.Executable
	})
}

func (p MigrationPlan) HasManualActions() bool {
	return lo.SomeBy(p.Actions, func(action MigrationAction) bool {
		return !action.Executable
	})
}
