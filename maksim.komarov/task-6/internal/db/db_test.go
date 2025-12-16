package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestDBService_GetNames_OK(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer sqlDB.Close()

	mockRows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).
		WillReturnRows(mockRows)

	service := New(sqlDB)

	got, err := service.GetNames()
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if strings.Join(got, ",") != "Alice,Bob" {
		t.Fatalf("unexpected names: %#v", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDBService_GetNames_QueryError(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer sqlDB.Close()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).
		WillReturnError(errors.New("db down"))

	service := New(sqlDB)

	_, e := service.GetNames()
	if e == nil || !strings.Contains(e.Error(), "db query:") {
		t.Fatalf("expected wrapped query error, got: %v", e)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDBService_GetNames_ScanError(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer sqlDB.Close()

	mockRows := sqlmock.NewRows([]string{"name"}).AddRow(nil)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).
		WillReturnRows(mockRows)

	service := New(sqlDB)

	_, e := service.GetNames()
	if e == nil || !strings.Contains(e.Error(), "rows scanning:") {
		t.Fatalf("expected wrapped scan error, got: %v", e)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDBService_GetNames_RowsErr(t *testing.T) {
	sqlDB := openErrAfterFirstDB(t)
	defer sqlDB.Close()

	service := New(sqlDB)

	_, e := service.GetNames()
	if e == nil || !strings.Contains(e.Error(), "rows error:") {
		t.Fatalf("expected wrapped rows error, got: %v", e)
	}
}

func TestDBService_GetUniqueNames_OK(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer sqlDB.Close()

	mockRows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Alice").
		AddRow("Bob")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).
		WillReturnRows(mockRows)

	service := New(sqlDB)

	got, err := service.GetUniqueNames()
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if strings.Join(got, ",") != "Alice,Alice,Bob" {
		t.Fatalf("unexpected values: %#v", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDBService_GetUniqueNames_QueryError(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer sqlDB.Close()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).
		WillReturnError(errors.New("db down"))

	service := New(sqlDB)

	_, e := service.GetUniqueNames()
	if e == nil || !strings.Contains(e.Error(), "db query:") {
		t.Fatalf("expected wrapped query error, got: %v", e)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDBService_GetUniqueNames_ScanError(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer sqlDB.Close()

	mockRows := sqlmock.NewRows([]string{"name"}).AddRow(nil)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).
		WillReturnRows(mockRows)

	service := New(sqlDB)

	_, e := service.GetUniqueNames()
	if e == nil || !strings.Contains(e.Error(), "rows scanning:") {
		t.Fatalf("expected wrapped scan error, got: %v", e)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDBService_GetUniqueNames_RowsErr(t *testing.T) {
	sqlDB := openErrAfterFirstDB(t)
	defer sqlDB.Close()

	service := New(sqlDB)

	_, e := service.GetUniqueNames()
	if e == nil || !strings.Contains(e.Error(), "rows error:") {
		t.Fatalf("expected wrapped rows error, got: %v", e)
	}
}

var regOnce sync.Once

func openErrAfterFirstDB(t *testing.T) *sql.DB {
	t.Helper()

	regOnce.Do(func() {
		sql.Register("errafterfirst", errAfterFirstDriver{})
	})

	db, err := sql.Open("errafterfirst", "")
	if err != nil {
		t.Fatalf("sql.Open: %v", err)
	}
	_ = db.PingContext(context.Background())
	return db
}

type errAfterFirstDriver struct{}

func (errAfterFirstDriver) Open(name string) (driver.Conn, error) {
	return &errAfterFirstConn{}, nil
}

type errAfterFirstConn struct{}

func (*errAfterFirstConn) Prepare(query string) (driver.Stmt, error) { return &noopStmt{}, nil }
func (*errAfterFirstConn) Close() error                              { return nil }
func (*errAfterFirstConn) Begin() (driver.Tx, error)                 { return &noopTx{}, nil }

func (*errAfterFirstConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	return &errAfterFirstRows{}, nil
}

func (*errAfterFirstConn) Query(query string, args []driver.Value) (driver.Rows, error) {
	return &errAfterFirstRows{}, nil
}

type errAfterFirstRows struct {
	n int
}

func (*errAfterFirstRows) Columns() []string { return []string{"name"} }
func (*errAfterFirstRows) Close() error      { return nil }

func (r *errAfterFirstRows) Next(dest []driver.Value) error {
	if r.n == 0 {
		dest[0] = "Alice"
		r.n++
		return nil
	}
	return errors.New("rows error from Next")
}

type noopStmt struct{}

func (*noopStmt) Close() error  { return nil }
func (*noopStmt) NumInput() int { return -1 }
func (*noopStmt) Exec(args []driver.Value) (driver.Result, error) {
	return nil, errors.New("not supported")
}
func (*noopStmt) Query(args []driver.Value) (driver.Rows, error) {
	return nil, errors.New("not supported")
}

type noopTx struct{}

func (*noopTx) Commit() error   { return nil }
func (*noopTx) Rollback() error { return nil }
