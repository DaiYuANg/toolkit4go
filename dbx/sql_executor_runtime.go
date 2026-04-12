package dbx

import (
	"context"
	"database/sql"
	"errors"

	"github.com/samber/oops"
)

func sessionExecutor(session Session) (*SQLExecutor, error) {
	if session == nil {
		return nil, oops.In("dbx").
			With("op", "sql_session").
			Wrapf(ErrNilDB, "validate session")
	}
	exec := session.SQL()
	if exec == nil {
		return nil, oops.In("dbx").
			With("op", "sql_session").
			Wrapf(ErrNilDB, "resolve sql executor")
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
	rows, _, err := queryStatementBoundRows(ctx, executor, statement, params)
	return rows, err
}

func queryStatementBoundRows(ctx context.Context, executor *SQLExecutor, statement SQLStatementSource, params any) (*sql.Rows, BoundQuery, error) {
	if executor == nil {
		return nil, BoundQuery{}, oops.In("dbx").
			With("op", "sql_query_rows", "statement", statementName(statement)).
			Wrapf(ErrNilDB, "validate sql executor")
	}
	bound, err := executor.Bind(statement, params)
	if err != nil {
		return nil, BoundQuery{}, oops.In("dbx").
			With("op", "sql_query_rows", "statement", statementName(statement)).
			Wrapf(err, "bind statement rows")
	}
	rows, err := executor.queryBound(ctx, bound)
	if err != nil {
		return nil, BoundQuery{}, oops.In("dbx").
			With("op", "sql_query_rows", "statement", statementName(statement)).
			Wrapf(err, "query statement rows")
	}
	return rows, bound, nil
}

func statementName(statement SQLStatementSource) string {
	if statement == nil {
		return ""
	}
	return statement.StatementName()
}
