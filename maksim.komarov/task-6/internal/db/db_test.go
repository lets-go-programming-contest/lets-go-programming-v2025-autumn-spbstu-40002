package db_test

import (
	"database/sql"
	"io"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	db "github.com/megurumacabre/task-6/internal/db"
)

func TestDBService_GetNames_OK(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}

	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).
		WillReturnRows(rows)

	service := db.New(sqlDB)

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
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}

	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).
		WillReturnError(sql.ErrConnDone)

	service := db.New(sqlDB)

	_, gotErr := service.GetNames()
	if gotErr == nil || !strings.Contains(gotErr.Error(), "db query:") {
		t.Fatalf("expected wrapped query error, got: %v", gotErr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDBService_GetNames_ScanError(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}

	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).
		WillReturnRows(rows)

	service := db.New(sqlDB)

	_, gotErr := service.GetNames()
	if gotErr == nil || !strings.Contains(gotErr.Error(), "rows scanning:") {
		t.Fatalf("expected wrapped scan error, got: %v", gotErr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDBService_GetNames_RowsErr(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}

	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob").
		RowError(1, io.ErrUnexpectedEOF)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).
		WillReturnRows(rows)

	service := db.New(sqlDB)

	_, gotErr := service.GetNames()
	if gotErr == nil || !strings.Contains(gotErr.Error(), "rows error:") {
		t.Fatalf("expected wrapped rows error, got: %v", gotErr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDBService_GetUniqueNames_OK(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}

	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Alice").
		AddRow("Bob")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).
		WillReturnRows(rows)

	service := db.New(sqlDB)

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
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}

	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).
		WillReturnError(sql.ErrConnDone)

	service := db.New(sqlDB)

	_, gotErr := service.GetUniqueNames()
	if gotErr == nil || !strings.Contains(gotErr.Error(), "db query:") {
		t.Fatalf("expected wrapped query error, got: %v", gotErr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDBService_GetUniqueNames_ScanError(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}

	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).
		WillReturnRows(rows)

	service := db.New(sqlDB)

	_, gotErr := service.GetUniqueNames()
	if gotErr == nil || !strings.Contains(gotErr.Error(), "rows scanning:") {
		t.Fatalf("expected wrapped scan error, got: %v", gotErr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDBService_GetUniqueNames_RowsErr(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}

	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob").
		RowError(1, io.ErrUnexpectedEOF)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).
		WillReturnRows(rows)

	service := db.New(sqlDB)

	_, gotErr := service.GetUniqueNames()
	if gotErr == nil || !strings.Contains(gotErr.Error(), "rows error:") {
		t.Fatalf("expected wrapped rows error, got: %v", gotErr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
