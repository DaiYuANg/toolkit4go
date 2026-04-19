package dbx_test

import (
	dbx "github.com/DaiYuANg/arcgo/dbx"
	codecx "github.com/DaiYuANg/arcgo/dbx/codec"
	"github.com/DaiYuANg/arcgo/dbx/idgen"
	"github.com/DaiYuANg/arcgo/dbx/paging"
	"github.com/DaiYuANg/arcgo/dbx/querydsl"
	relationx "github.com/DaiYuANg/arcgo/dbx/relation"
	schemax "github.com/DaiYuANg/arcgo/dbx/schema"
	schemamigrate "github.com/DaiYuANg/arcgo/dbx/schemamigrate"
	"github.com/DaiYuANg/arcgo/dbx/sqlexec"
	"github.com/DaiYuANg/arcgo/dbx/sqlstmt"
)

type Aggregate[T any] = querydsl.Aggregate[T]
type Assignment = querydsl.Assignment
type AtlasCompiledSchemaTestView = schemamigrate.AtlasCompiledSchemaTestView
type BelongsTo[E any, T any] = relationx.BelongsTo[E, T]
type Bound = sqlstmt.Bound
type CaseBuilder[T any] = querydsl.CaseBuilder[T]
type Check[E any] = dbx.Check[E]
type CheckMeta = schemax.CheckMeta
type CheckState = schemax.CheckState
type Codec = codecx.Codec
type Column[E any, T any] = dbx.Column[E, T]
type ColumnMeta = schemax.ColumnMeta
type ColumnState = schemax.ColumnState
type CompositeKey[E any] = dbx.CompositeKey[E]
type Cursor[T any] = dbx.Cursor[T]
type DB = dbx.DB
type DeleteQuery = querydsl.DeleteQuery
type Executor = dbx.Executor
type ForeignKeyMeta = schemax.ForeignKeyMeta
type ForeignKeyState = schemax.ForeignKeyState
type HasMany[E any, T any] = relationx.HasMany[E, T]
type HasOne[E any, T any] = relationx.HasOne[E, T]
type Hook = dbx.Hook
type HookEvent = dbx.HookEvent
type HookFuncs = dbx.HookFuncs
type IDGenerator = idgen.Generator
type IDColumn[E any, T any, M dbx.IDMarker] = dbx.IDColumn[E, T, M]
type IDKSUID = dbx.IDKSUID
type IDSnowflake = dbx.IDSnowflake
type IDULID = dbx.IDULID
type IDUUID = dbx.IDUUID
type IDUUIDv4 = dbx.IDUUIDv4
type IDUUIDv7 = dbx.IDUUIDv7
type Index[E any] = dbx.Index[E]
type IndexMeta = schemax.IndexMeta
type IndexState = schemax.IndexState
type InsertQuery = querydsl.InsertQuery
type MappedField = dbx.MappedField
type ManyToMany[E any, T any] = relationx.ManyToMany[E, T]
type Mapper[E any] = dbx.Mapper[E]
type MapperOption = dbx.MapperOption
type MigrationAction = schemax.MigrationAction
type MigrationPlan = schemax.MigrationPlan
type NodeIDOutOfRangeError = idgen.NodeIDOutOfRangeError
type Operation = dbx.Operation
type Option = dbx.Option
type PageRequest = paging.Request
type PageResult[E any] = paging.Result[E]
type Predicate = querydsl.Predicate
type PrimaryKeyMeta = schemax.PrimaryKeyMeta
type PrimaryKeyState = schemax.PrimaryKeyState
type PrimaryKeyUnmappedError = dbx.PrimaryKeyUnmappedError
type QueryBuilder = querydsl.Builder
type ReferentialAction = schemax.ReferentialAction
type RelationKind = schemax.RelationKind
type Row = dbx.Row
type RowsScanner[E any] = dbx.RowsScanner[E]
type Schema[E any] = dbx.Schema[E]
type SchemaDriftError = schemax.SchemaDriftError
type SchemaResource = schemamigrate.Resource
type SchemaSource[E any] = dbx.SchemaSource[E]
type SelectItem = querydsl.SelectItem
type SelectQuery = querydsl.SelectQuery
type Session = dbx.Session
type SQLExecutor = sqlexec.Executor
type Statement = sqlstmt.Statement
type StatementSource = sqlstmt.Source
type StructMapper[E any] = dbx.StructMapper[E]
type Table = querydsl.Table
type TableSource = querydsl.TableSource
type TableSpec = schemax.TableSpec
type TableState = schemax.TableState
type UnknownCodecError = codecx.UnknownError
type UnmappedColumnError = dbx.UnmappedColumnError
type ValidationBackend = schemax.ValidationBackend
type ValidationReport = schemax.ValidationReport

const (
	DefaultNodeID              = idgen.DefaultNodeID
	DefaultPage                = paging.DefaultPage
	DefaultPageSize            = paging.DefaultPageSize
	DefaultUUIDVersion         = idgen.DefaultUUIDVersion
	IDStrategyDBAuto           = idgen.StrategyDBAuto
	IDStrategySnowflake        = idgen.StrategySnowflake
	IDStrategyUUID             = idgen.StrategyUUID
	MigrationActionCreateIndex = schemax.MigrationActionCreateIndex
	MigrationActionCreateTable = schemax.MigrationActionCreateTable
	MigrationActionManual      = schemax.MigrationActionManual
	OperationAutoMigrate       = dbx.OperationAutoMigrate
	OperationBeginTx           = dbx.OperationBeginTx
	OperationCommitTx          = dbx.OperationCommitTx
	OperationExec              = dbx.OperationExec
	OperationQuery             = dbx.OperationQuery
	OperationQueryRow          = dbx.OperationQueryRow
	OperationRollbackTx        = dbx.OperationRollbackTx
	OperationValidate          = dbx.OperationValidate
	ReferentialCascade         = schemax.ReferentialCascade
	RelationBelongsTo          = schemax.RelationBelongsTo
	RelationManyToMany         = schemax.RelationManyToMany
	ValidationBackendLegacy    = schemax.ValidationBackendLegacy
)
