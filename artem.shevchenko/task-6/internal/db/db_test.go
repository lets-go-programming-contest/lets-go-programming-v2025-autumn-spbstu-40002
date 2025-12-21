package db_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/slendycs/go-lab-6/internal/db"
	"github.com/stretchr/testify/require"
)

var (
	errConnectionFailed    = errors.New("connection failed")
	errRowsIterationFailed = errors.New("rows iteration failed")
	errUniqueQueryFailed   = errors.New("unique query failed")
	errIterationError      = errors.New("iteration error")
	errRowErrorTest        = errors.New("row error")
	errCustomQueryError    = errors.New("custom query error")
	errCloseErrorTest      = errors.New("close error")
)

func TestGetNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		mockRows     *sqlmock.Rows
		mockError    error
		expectedErr  string
		expectedData []string
	}{
		{
			name:         "success - multiple names",
			mockRows:     sqlmock.NewRows([]string{"name"}).AddRow("Ivan").AddRow("Gena228").AddRow("Petr"),
			expectedData: []string{"Ivan", "Gena228", "Petr"},
		},
		{
			name:         "success - empty result",
			mockRows:     sqlmock.NewRows([]string{"name"}),
			expectedData: []string(nil),
		},
		{
			name:        "query error",
			mockError:   errConnectionFailed,
			expectedErr: "db query: connection failed",
		},
		{
			name:         "rows error after iteration",
			mockRows:     sqlmock.NewRows([]string{"name"}).AddRow("Ivan").CloseError(errRowsIterationFailed),
			expectedErr:  "rows error: rows iteration failed",
			expectedData: []string{"Ivan"},
		},
		{
			name:         "scan error - null value",
			mockRows:     sqlmock.NewRows([]string{"name"}).AddRow(nil),
			expectedErr:  "rows scanning:",
			expectedData: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			expect := mock.ExpectQuery("SELECT name FROM users")

			if tc.mockError != nil {
				expect.WillReturnError(tc.mockError)
			} else if tc.mockRows != nil {
				expect.WillReturnRows(tc.mockRows)
			}

			dbService := db.New(mockDB)
			names, err := dbService.GetNames()

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
				require.Nil(t, names)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedData, names)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUniqueNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		mockRows     *sqlmock.Rows
		mockError    error
		expectedErr  string
		expectedData []string
	}{
		{
			name:         "success - unique names",
			mockRows:     sqlmock.NewRows([]string{"name"}).AddRow("Ivan").AddRow("Petr").AddRow("Ivan").AddRow("Sergey"),
			expectedData: []string{"Ivan", "Petr", "Ivan", "Sergey"},
		},
		{
			name:         "success - single name",
			mockRows:     sqlmock.NewRows([]string{"name"}).AddRow("Ivan"),
			expectedData: []string{"Ivan"},
		},
		{
			name:        "query error",
			mockError:   errUniqueQueryFailed,
			expectedErr: "db query: unique query failed",
		},
		{
			name:         "rows error after iteration",
			mockRows:     sqlmock.NewRows([]string{"name"}).AddRow("Ivan").AddRow("Petr").CloseError(errIterationError),
			expectedErr:  "rows error: iteration error",
			expectedData: []string{"Ivan", "Petr"},
		},
		{
			name:         "success - no rows",
			mockRows:     sqlmock.NewRows([]string{"name"}),
			expectedData: []string(nil),
		},
		{
			name:         "scan error - null value",
			mockRows:     sqlmock.NewRows([]string{"name"}).AddRow(nil),
			expectedErr:  "rows scanning:",
			expectedData: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			expect := mock.ExpectQuery("SELECT DISTINCT name FROM users")

			if tc.mockError != nil {
				expect.WillReturnError(tc.mockError)
			} else if tc.mockRows != nil {
				expect.WillReturnRows(tc.mockRows)
			}

			dbService := db.New(mockDB)
			names, err := dbService.GetUniqueNames()

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
				require.Nil(t, names)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedData, names)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestNoExpectations(t *testing.T) {
	t.Parallel()

	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	dbService := db.New(mockDB)
	names, err := dbService.GetNames()

	require.Error(t, err)
	require.Contains(t, err.Error(), "all expectations were already fulfilled")
	require.Nil(t, names)
}

func TestDBIntegration(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Ivan").AddRow("Petr"))

	dbService := db.New(mockDB)
	names, err := dbService.GetNames()
	require.NoError(t, err)
	require.Equal(t, []string{"Ivan", "Petr"}, names)

	mock.ExpectQuery("SELECT DISTINCT name FROM users").
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Ivan").AddRow("Petr"))

	uniqueNames, err := dbService.GetUniqueNames()
	require.NoError(t, err)
	require.Equal(t, []string{"Ivan", "Petr"}, uniqueNames)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRowsAreClosed(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"name"}).AddRow("Ivan")
	rows.CloseError(nil)

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	dbService := db.New(mockDB)
	names, err := dbService.GetNames()

	require.NoError(t, err)
	require.Equal(t, []string{"Ivan"}, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNamesRowError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Ivan").
		AddRow("Petr").
		RowError(1, errRowErrorTest)

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	dbService := db.New(mockDB)
	names, err := dbService.GetNames()

	require.Error(t, err)
	require.Contains(t, err.Error(), "rows error: row error")
	require.Nil(t, names)
}

type errorDatabase struct {
	err error
}

func (e *errorDatabase) Query(query string, args ...any) (*sql.Rows, error) {
	return nil, e.err
}

func TestGetNamesWithErrorDatabase(t *testing.T) {
	t.Parallel()

	mockError := errCustomQueryError

	dbService := db.New(&errorDatabase{err: mockError})
	names, err := dbService.GetNames()

	require.Error(t, err)
	require.Contains(t, err.Error(), "db query: custom query error")
	require.Nil(t, names)
}

func TestGetUniqueNamesWithErrorDatabase(t *testing.T) {
	t.Parallel()

	mockError := errCustomQueryError

	dbService := db.New(&errorDatabase{err: mockError})
	names, err := dbService.GetUniqueNames()

	require.Error(t, err)
	require.Contains(t, err.Error(), "db query: custom query error")
	require.Nil(t, names)
}

func TestGetNamesRowsErrAfterCompleteScan(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Ivan").
		AddRow("Petr").
		CloseError(errCloseErrorTest)

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	dbService := db.New(mockDB)
	names, err := dbService.GetNames()

	require.Error(t, err)
	require.Contains(t, err.Error(), "rows error: close error")
	require.Nil(t, names)
}
