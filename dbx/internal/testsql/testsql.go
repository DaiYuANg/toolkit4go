package testsql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"reflect"
	"sync"
	"sync/atomic"
)

const DriverName = "dbx-internal-testsql"

type QueryPlan struct {
	SQL     string
	Args    []driver.Value
	Columns []string
	Rows    [][]driver.Value
	Err     error
}

type ExecPlan struct {
	SQL          string
	Args         []driver.Value
	RowsAffected int64
	LastInsertID int64
	Err          error
}

type Plan struct {
	Queries []QueryPlan
	Execs   []ExecPlan
}

type Call struct {
	SQL  string
	Args []driver.Value
}

type Recorder struct {
	mu      sync.Mutex
	Queries []Call
	Execs   []Call
}

type script struct {
	mu       sync.Mutex
	queries  []QueryPlan
	execs    []ExecPlan
	recorder *Recorder
}

type testDriver struct{}

type testConn struct {
	script *script
}

type testTx struct{}

type testRows struct {
	columns []string
	rows    [][]driver.Value
	index   int
}

type testResult struct {
	rowsAffected int64
	lastInsertID int64
}

var (
	registerOnce sync.Once
	scripts      sync.Map
	sequence     uint64
)

func Open(plan Plan) (*sql.DB, *Recorder, func(), error) {
	registerOnce.Do(func() {
		sql.Register(DriverName, testDriver{})
	})

	key := fmt.Sprintf("plan-%d", atomic.AddUint64(&sequence, 1))
	recorder := &Recorder{
		Queries: make([]Call, 0, len(plan.Queries)),
		Execs:   make([]Call, 0, len(plan.Execs)),
	}
	scripts.Store(key, &script{
		queries:  append([]QueryPlan(nil), plan.Queries...),
		execs:    append([]ExecPlan(nil), plan.Execs...),
		recorder: recorder,
	})

	db, err := sql.Open(DriverName, key)
	if err != nil {
		scripts.Delete(key)
		return nil, nil, nil, err
	}

	cleanup := func() {
		_ = db.Close()
		scripts.Delete(key)
	}
	return db, recorder, cleanup, nil
}

func (testDriver) Open(name string) (driver.Conn, error) {
	value, ok := scripts.Load(name)
	if !ok {
		return nil, fmt.Errorf("testsql: unknown script %q", name)
	}
	return &testConn{script: value.(*script)}, nil
}

func (c *testConn) Prepare(string) (driver.Stmt, error) {
	return nil, fmt.Errorf("testsql: Prepare is not supported")
}

func (c *testConn) Close() error {
	return nil
}

func (c *testConn) Begin() (driver.Tx, error) {
	return testTx{}, nil
}

func (c *testConn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	return testTx{}, nil
}

func (c *testConn) QueryContext(_ context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	item, err := c.script.popQuery(query, args)
	if err != nil {
		return nil, err
	}
	if item.Err != nil {
		return nil, item.Err
	}
	return &testRows{
		columns: append([]string(nil), item.Columns...),
		rows:    cloneRows(item.Rows),
	}, nil
}

func (c *testConn) ExecContext(_ context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	item, err := c.script.popExec(query, args)
	if err != nil {
		return nil, err
	}
	if item.Err != nil {
		return nil, item.Err
	}
	return testResult{
		rowsAffected: item.RowsAffected,
		lastInsertID: item.LastInsertID,
	}, nil
}

func (s *script) popQuery(query string, args []driver.NamedValue) (QueryPlan, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	call := Call{SQL: query, Args: namedValues(args)}
	s.recorder.Queries = append(s.recorder.Queries, call)
	if len(s.queries) == 0 {
		return QueryPlan{}, fmt.Errorf("testsql: unexpected query %q", query)
	}

	item := s.queries[0]
	s.queries = s.queries[1:]
	if err := compareCall(item.SQL, item.Args, call); err != nil {
		return QueryPlan{}, err
	}
	return item, nil
}

func (s *script) popExec(query string, args []driver.NamedValue) (ExecPlan, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	call := Call{SQL: query, Args: namedValues(args)}
	s.recorder.Execs = append(s.recorder.Execs, call)
	if len(s.execs) == 0 {
		return ExecPlan{}, fmt.Errorf("testsql: unexpected exec %q", query)
	}

	item := s.execs[0]
	s.execs = s.execs[1:]
	if err := compareCall(item.SQL, item.Args, call); err != nil {
		return ExecPlan{}, err
	}
	return item, nil
}

func namedValues(args []driver.NamedValue) []driver.Value {
	values := make([]driver.Value, len(args))
	for i, arg := range args {
		values[i] = arg.Value
	}
	return values
}

func compareCall(expectedSQL string, expectedArgs []driver.Value, actual Call) error {
	if expectedSQL != "" && expectedSQL != actual.SQL {
		return fmt.Errorf("testsql: unexpected sql, want %q got %q", expectedSQL, actual.SQL)
	}
	if expectedArgs != nil && !reflect.DeepEqual(expectedArgs, actual.Args) {
		return fmt.Errorf("testsql: unexpected args, want %#v got %#v", expectedArgs, actual.Args)
	}
	return nil
}

func cloneRows(rows [][]driver.Value) [][]driver.Value {
	items := make([][]driver.Value, len(rows))
	for i, row := range rows {
		items[i] = append([]driver.Value(nil), row...)
	}
	return items
}

func (r *testRows) Columns() []string {
	return append([]string(nil), r.columns...)
}

func (r *testRows) Close() error {
	return nil
}

func (r *testRows) Next(dest []driver.Value) error {
	if r.index >= len(r.rows) {
		return io.EOF
	}
	row := r.rows[r.index]
	r.index++
	copy(dest, row)
	return nil
}

func (r testResult) LastInsertId() (int64, error) {
	return r.lastInsertID, nil
}

func (r testResult) RowsAffected() (int64, error) {
	return r.rowsAffected, nil
}

func (testTx) Commit() error {
	return nil
}

func (testTx) Rollback() error {
	return nil
}
