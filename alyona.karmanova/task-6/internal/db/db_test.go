package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	myDB "github.com/HuaChenju/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

var testTable = []struct { //nolint:gochecknoglobals
	names       []string
	errWrap     string
	errExpected error
}{
	{
		names: []string{"HuaCheng", "XieLian"},
	},
	{
		names:       nil,
		errWrap:     "db query: ",
		errExpected: ErrExpected,
	},
}

var (
	ErrExpected       = errors.New("errExpected")
	QueryName         = "SELECT name FROM users"          //nolint:gochecknoglobals
	QueryDistinctName = "SELECT DISTINCT name FROM users" //nolint:gochecknoglobals
)

func TestNew(t *testing.T) {
	t.Parallel()

	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := myDB.New(mockDB)

	require.Equal(t, mockDB, service.DB)
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected jsondata", err)
	}

	defer mockDB.Close()

	service := myDB.DBService{DB: mockDB}

	for i, row := range testTable {
		mock.ExpectQuery(QueryName).WillReturnRows(helperMockDBRows(t, row.names)).WillReturnError(row.errExpected)
		names, err := service.GetNames()

		if row.errExpected != nil {
			require.ErrorIs(t, err, row.errExpected, "row: %d, expected error:%w, actual error: %w", i, row.errExpected, err)
			require.Nil(t, names, "row: %d, names must be nil", i)

			continue
		}

		require.NoError(t, err, "row: %d, error must be nil", i)
		require.Equal(t, row.names, names, "row: %d, expected names: %s, actual names: %s", i,
			row.names, names)
	}

	mock.ExpectQuery(QueryName).WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))

	names, err := service.GetNames()

	require.ErrorContains(t, err, "rows scanning")
	require.Nil(t, names, "must be nil")

	mock.ExpectQuery(QueryName).WillReturnRows(sqlmock.NewRows([]string{"name"}).
		AddRow("HuaCheng").AddRow("XieLian").RowError(1, ErrExpected))

	names, err = service.GetNames()
	require.ErrorContains(t, err, "rows error")
	require.Nil(t, names, "must be nil")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetDistinctNames(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected jsondata", err)
	}

	defer mockDB.Close()

	service := myDB.DBService{DB: mockDB}

	for i, row := range testTable {
		mock.ExpectQuery(QueryDistinctName).WillReturnRows(helperMockDBRows(t, row.names)).WillReturnError(row.errExpected)
		names, err := service.GetUniqueNames()

		if row.errExpected != nil {
			require.ErrorIs(t, err, row.errExpected, "row: %d, expected error:%w, actual error: %w", i, row.errExpected, err)
			require.Nil(t, names, "row: %d, names must be nil", i)

			continue
		}

		require.NoError(t, err, "row: %d, error must be nil", i)
		require.Equal(t, row.names, names, "row: %d, expected names: %s, actual names: %s", i,
			row.names, names)
	}

	mock.ExpectQuery(QueryDistinctName).WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))

	names, err := service.GetUniqueNames()

	require.ErrorContains(t, err, "rows scanning")
	require.Nil(t, names, "must be nil")

	mock.ExpectQuery(QueryDistinctName).WillReturnRows(sqlmock.NewRows([]string{"name"}).
		AddRow("HuaCheng").AddRow("XieLian").RowError(1, ErrExpected))

	names, err = service.GetUniqueNames()

	require.ErrorContains(t, err, "rows error")
	require.Nil(t, names, "must be nil")
	require.NoError(t, mock.ExpectationsWereMet())
}

func helperMockDBRows(t *testing.T, names []string) *sqlmock.Rows {
	t.Helper()

	rows := sqlmock.NewRows([]string{"name"})

	for _, name := range names {
		rows = rows.AddRow(name)
	}

	return rows
}
