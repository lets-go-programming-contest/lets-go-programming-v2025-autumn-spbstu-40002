package db_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	database "github.com/ControlShiftEscape/task-6/internal/db"
)

var (
	errScanFailed     = errors.New("scan failed")
	errIteratorFailed = errors.New("iterator error")
	errCloseFailed    = errors.New("close failed")
	errRowsIterator   = errors.New("rows iterator failed")
)

type TestCase struct {
	name      string
	setupMock func(sqlmock.Sqlmock)
	wantNames []string
	wantErr   bool
}

func TestNew(t *testing.T) {
	t.Parallel()

	db, _, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	service := database.New(db)

	require.NotNil(t, service)
	require.Equal(t, db, service.DB)
}

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) { //nolint:ireturn
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	return db, mock
}

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	tests := []TestCase{
		{
			name: "success",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alexander").AddRow("Petr").AddRow("Feodosiy")
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantNames: []string{"Alexander", "Petr", "Feodosiy"},
			wantErr:   false,
		},
		{
			name: "db_error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnError(sql.ErrConnDone)
			},
			wantNames: nil,
			wantErr:   true,
		},
		{
			name: "scan_error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					RowError(0, errScanFailed)
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantNames: nil,
			wantErr:   true,
		},
		{
			name: "rows_error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alexander").
					RowError(1, errIteratorFailed)
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantNames: nil,
			wantErr:   true,
		},
		{
			name: "empty_result",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantNames: []string{},
			wantErr:   false,
		},
		{
			name: "rows_close_error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alexander").
					CloseError(errCloseFailed)
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantNames: nil,
			wantErr:   true,
		},
		{
			name: "scan_nil_value",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow(nil)
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantNames: nil,
			wantErr:   true,
		},
		{
			name: "rows_err_after_successful_query",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alexander").
					AddRow("Petr")
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnRows(rows).
					WillReturnError(errRowsIterator)
			},
			wantNames: nil,
			wantErr:   true,
		},
	}

	runTests(t, tests, func(s *database.DBService) ([]string, error) {
		return s.GetNames()
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	tests := []TestCase{
		{
			name: "success",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alexander").
					AddRow("Petr").
					AddRow("Alexander").
					AddRow("Feodosiy")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantNames: []string{"Alexander", "Petr", "Feodosiy"},
			wantErr:   false,
		},
		{
			name: "db_error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnError(sql.ErrConnDone)
			},
			wantNames: nil,
			wantErr:   true,
		},
		{
			name: "scan_error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					RowError(0, errScanFailed)
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantNames: nil,
			wantErr:   true,
		},
		{
			name: "rows_error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alexander").
					RowError(1, errRowsIterator)
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantNames: nil,
			wantErr:   true,
		},
		{
			name: "empty_result",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantNames: []string{},
			wantErr:   false,
		},
		{
			name: "rows_close_error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alexander").
					CloseError(errCloseFailed)
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantNames: nil,
			wantErr:   true,
		},
		{
			name: "scan_nil_value",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow(nil)
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantNames: nil,
			wantErr:   true,
		},
		{
			name: "rows_err_after_successful_query",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alexander").
					AddRow("Petr")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnRows(rows).
					WillReturnError(errRowsIterator)
			},
			wantNames: nil,
			wantErr:   true,
		},
	}

	runTests(t, tests, func(s *database.DBService) ([]string, error) {
		return s.GetUniqueNames()
	})
}

func runTests(t *testing.T, tests []TestCase, testFunc func(*database.DBService) ([]string, error)) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock := setupMockDB(t)

			defer db.Close()

			tt.setupMock(mock)

			service := database.DBService{DB: db}
			names, err := testFunc(&service)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, names)
			} else {
				require.NoError(t, err)
				require.NotNil(t, names)
				assert.Equal(t, tt.wantNames, names)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
