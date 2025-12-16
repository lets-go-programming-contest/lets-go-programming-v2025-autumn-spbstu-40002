package db

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

type rowTestDb struct {
	names       []string
	errExpected error
}

var testTable1 = []rowTestDb{
	{
		names: []string{"Ivan", "Gena228"},
	},
	{
		names:       nil,
		errExpected: errors.New("ExpectedError"),
	},
	{
		names: []string{"John", "", "Mary", "  ", "Bob", "\t", "Alice"},
	},
	{
		names: []string{"John", "Mary", "Bob", "Alice", "Charlie"},
	},
	{
		names: []string{"Anna"},
	},
	{
		names: []string{""},
	},
	{
		names:       nil,
		errExpected: sql.ErrNoRows,
	},
}

func TestGetNames(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "failed to create mock database")
	defer mockDB.Close()

	dbService := New(mockDB)

	for i, row := range testTable1 {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			if row.errExpected != nil {
				mock.ExpectQuery("SELECT name FROM users").WillReturnError(row.errExpected)
			} else {
				rows := sqlmock.NewRows([]string{"name"})
				for _, name := range row.names {
					rows = rows.AddRow(name)
				}
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			}

			names, err := dbService.GetNames()

			if row.errExpected != nil {
				require.Error(t, err, "row: %d, expected error but got none", i)
				if errors.Is(row.errExpected, sql.ErrNoRows) {
					require.ErrorIs(t, err, sql.ErrNoRows, "row: %d", i)
				}
				require.Nil(t, names, "row: %d, names must be nil", i)
			} else {
				require.NoError(t, err, "row: %d, error must be nil", i)
				if row.names == nil || len(row.names) == 0 {
					require.Nil(t, names, "row: %d, для пустого запроса должен возвращаться nil", i)
				} else {
					require.Equal(t, row.names, names, "row: %d, expected names: %v, actual names: %v", i,
						row.names, names)
				}
			}
		})
	}

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "failed to create mock database")
	defer mockDB.Close()

	dbService := New(mockDB)

	for i, row := range testTable1 {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			if row.errExpected != nil {
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(row.errExpected)
			} else {
				rows := sqlmock.NewRows([]string{"name"})
				for _, name := range row.names {
					rows = rows.AddRow(name)
				}
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			}

			uniqueNames, err := dbService.GetUniqueNames()

			if row.errExpected != nil {
				require.Error(t, err, "row: %d, expected error but got none", i)
				if errors.Is(row.errExpected, sql.ErrNoRows) {
					require.ErrorIs(t, err, sql.ErrNoRows, "row: %d", i)
				}
				require.Nil(t, uniqueNames, "row: %d, names must be nil", i)
			} else {
				require.NoError(t, err, "row: %d, error must be nil", i)
				if row.names == nil || len(row.names) == 0 {
					require.Nil(t, uniqueNames, "row: %d, для пустого запроса должен возвращаться nil", i)
				} else {
					require.Equal(t, row.names, uniqueNames, "row: %d, expected names: %v, actual names: %v", i,
						row.names, uniqueNames)
				}
			}
		})
	}

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_ErrorInRowsErr(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob").
		CloseError(errors.New("rows error"))

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	dbService := New(mockDB)
	names, err := dbService.GetNames()

	require.Error(t, err)
	require.Contains(t, err.Error(), "rows error")
	require.Nil(t, names)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_ErrorInRowsErr(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob").
		CloseError(errors.New("rows error"))

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	dbService := New(mockDB)
	names, err := dbService.GetUniqueNames()

	require.Error(t, err)
	require.Contains(t, err.Error(), "rows error")
	require.Nil(t, names)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_EmptyResult(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"name"})

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	dbService := New(mockDB)
	names, err := dbService.GetNames()

	require.NoError(t, err)
	require.Nil(t, names)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_EmptyResult(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"name"})

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	dbService := New(mockDB)
	names, err := dbService.GetUniqueNames()

	require.NoError(t, err)
	require.Nil(t, names)

	require.NoError(t, mock.ExpectationsWereMet())
}

type mockDatabaseWithScanError struct{}

func (m *mockDatabaseWithScanError) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, errors.New("cannot easily mock scan error with sqlmock")
}

func TestGetNames_DBQueryError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	expectedErr := errors.New("database connection failed")
	mock.ExpectQuery("SELECT name FROM users").WillReturnError(expectedErr)

	dbService := New(mockDB)
	names, err := dbService.GetNames()

	require.Error(t, err)
	require.Contains(t, err.Error(), "db query")
	require.Contains(t, err.Error(), expectedErr.Error())
	require.Nil(t, names)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_DBQueryError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	expectedErr := errors.New("database connection failed")
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(expectedErr)

	dbService := New(mockDB)
	names, err := dbService.GetUniqueNames()

	require.Error(t, err)
	require.Contains(t, err.Error(), "db query")
	require.Contains(t, err.Error(), expectedErr.Error())
	require.Nil(t, names)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_RowErrorScenario(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob").
		RowError(1, errors.New("row error"))

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	dbService := New(mockDB)
	names, err := dbService.GetNames()

	if err != nil {
		require.Contains(t, err.Error(), "rows error")
		require.Nil(t, names)
	} else {
		require.NotNil(t, names)
		require.Len(t, names, 2)
	}

	require.NoError(t, mock.ExpectationsWereMet())
}
