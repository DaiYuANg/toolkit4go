package dbx

import (
	"context"
	"database/sql"
	"errors"

	"github.com/DaiYuANg/arcgo/collectionx"
	"github.com/samber/oops"
	scanlib "github.com/stephenafamo/scan"
)

func scanSQLListRowsWithCapacity[E any](session Session, rows *sql.Rows, bound BoundQuery, mapper CapacityHintScanner[E]) (collectionx.List[E], error) {
	logRuntimeNode(session, "sql.list.scan_with_capacity", "statement", bound.Name, "capacity_hint", bound.CapacityHint)
	items, scanErr := mapper.ScanRowsWithCapacity(rows, bound.CapacityHint)
	scanErr = errors.Join(wrapDBError("scan statement rows with capacity", scanErr), rowsIterError(rows))
	closeErr := closeRows(rows)
	if scanErr != nil {
		logRuntimeNode(session, "sql.list.error", "stage", "scan_rows_with_capacity", "error", scanErr)
		return nil, errors.Join(scanErr, closeErr)
	}
	if closeErr != nil {
		logRuntimeNode(session, "sql.list.error", "stage", "close_rows", "error", closeErr)
		return nil, closeErr
	}
	logRuntimeNode(session, "sql.list.done", "items", items.Len())
	return items, nil
}

func sqlScalar[T any](ctx context.Context, session Session, statement SQLStatementSource, params any) (T, bool, error) {
	exec, err := sessionExecutor(session)
	if err != nil {
		var zero T
		return zero, false, err
	}
	rows, err := queryStatementRows(ctx, exec, statement, params)
	if err != nil {
		var zero T
		return zero, false, err
	}

	value, err := scanlib.OneFromRows[T](ctx, scanlib.SingleColumnMapper[T], rows)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			closeErr := closeRows(rows)
			var zero T
			return zero, false, closeErr
		}
		closeErr := closeRows(rows)
		var zero T
		return zero, false, errors.Join(wrapDBError("scan scalar row", err), closeErr)
	}

	if rows.Next() {
		closeErr := closeRows(rows)
		var zero T
		return zero, false, errors.Join(
			oops.In("dbx").
				With("op", "sql_scalar", "statement", statementName(statement)).
				Wrapf(ErrTooManyRows, "sql scalar returned too many rows"),
			closeErr,
		)
	}
	if err := rowsIterError(rows); err != nil {
		closeErr := closeRows(rows)
		var zero T
		return zero, false, errors.Join(err, closeErr)
	}
	closeErr := closeRows(rows)
	if closeErr != nil {
		var zero T
		return zero, false, closeErr
	}
	return value, true, nil
}
