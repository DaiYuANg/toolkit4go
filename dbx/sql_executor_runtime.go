package dbx

import (
	"context"
	"database/sql"
	"errors"
)

func sessionExecutor(session Session) (*SQLExecutor, error) {
	if session == nil {
		return nil, ErrNilDB
	}
	exec := session.SQL()
	if exec == nil {
		return nil, ErrNilDB
	}
	return exec, nil
}

func scanStatementOne[E any](ctx context.Context, exec *SQLExecutor, statement SQLStatementSource, params any, mapper oneRowScanner[E]) (E, bool, error) {
	rows, err := queryStatementRows(ctx, exec, statement, params)
	if err != nil {
		var zero E
		return zero, false, err
	}
	value, found, scanErr := mapper.scanOneRows(ctx, rows)
	closeErr := closeRows(rows)
	if scanErr != nil {
		var zero E
		return zero, false, errors.Join(scanErr, closeErr)
	}
	if closeErr != nil {
		var zero E
		return zero, false, closeErr
	}
	return value, found, nil
}

func queryStatementRows(ctx context.Context, executor *SQLExecutor, statement SQLStatementSource, params any) (*sql.Rows, error) {
	if executor == nil {
		return nil, ErrNilDB
	}
	rows, err := executor.Query(ctx, statement, params)
	return rows, wrapDBError("query statement rows", err)
}
