package dbx

import (
	"database/sql"
	schemax "github.com/DaiYuANg/arcgo/dbx/schema"

	"github.com/DaiYuANg/arcgo/collectionx"
)

func TableSpecForTest(schema schemax.Resource) schemax.TableSpec {
	return schema.Spec()
}

func IndexesForTest(schema schemax.Resource) collectionx.List[schemax.IndexMeta] {
	return schema.Spec().Indexes.Clone()
}

func InferTypeNameForTest(column schemax.ColumnMeta) string {
	return schemax.InferTypeName(column)
}

func ErrorRowForTest(err error) *Row {
	return errorRow(err)
}

func CloseRowsForTest(rows *sql.Rows) error {
	return closeRows(rows)
}

func RowsIterErrorForTest(rows *sql.Rows) error {
	return rowsIterError(rows)
}

func StructMapperScanPlanForTest[E any](mapper StructMapper[E], columns []string) error {
	_, err := mapper.scanPlan(columns)
	return err
}

func ClonePrimaryKeyMetaForTest(meta schemax.PrimaryKeyMeta) schemax.PrimaryKeyMeta {
	return clonePrimaryKeyMeta(meta)
}

func ClonePrimaryKeyStateForTest(state schemax.PrimaryKeyState) schemax.PrimaryKeyState {
	return clonePrimaryKeyState(state)
}
