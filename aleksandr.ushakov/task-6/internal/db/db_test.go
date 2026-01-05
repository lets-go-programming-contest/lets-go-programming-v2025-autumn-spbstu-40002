package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	myDb "github.com/rachguta/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

var (
	errExpectedTest = errors.New("ExpectedError")
	errNetwork      = errors.New("network problem")
)

type rowTestDB struct {
	names       []string
	errExpected error
}

var testTable = []rowTestDB{ //nolint:gochecknoglobals
	{
		names: []string{"Ivan", "Gena228"},
	},
	{
		names:       nil,
		errExpected: errExpectedTest,
	},
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dbService := myDb.DBService{DB: mockDB}

	for i, row := range testTable {
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(mockDBRows(row.names)).
			WillReturnError(row.errExpected)

		names, err := dbService.GetNames()

		if row.errExpected != nil {
			require.ErrorIs(t, err, row.errExpected, "row: %d, expected error: %w, actual error: %w", i,
				row.errExpected, err)
			require.Nil(t, names, "row: %d, names must be nil", i)

			continue
		}

		require.NoError(t, err, "row: %d, error must be nil", i)
		require.Equal(t, row.names, names, "row: %d, expected names: %s, actual names: %s",
			i, row.names, names)
	}
}

func TestGetNames_ScanError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dbService := myDb.DBService{DB: mockDB}
	// setting up the mock so that the query returns two columns
	rows := sqlmock.NewRows([]string{"name", "ages"}).AddRow("Ivan", "18")
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	names, err := dbService.GetNames()
	require.Error(t, err, "error must not be nil")
	require.Containsf(t, err.Error(), "rows scanning: ", "error: %s", err)
	require.Nil(t, names, "names must be nil")
}

func TestGetNames_RowError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dbService := myDb.DBService{DB: mockDB}
	// simulate a connection error when reading a row
	rows := sqlmock.NewRows([]string{"name"}).AddRow("Ivan").
		AddRow("Gena228").RowError(1, errNetwork)
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	names, err := dbService.GetNames()

	require.Error(t, err, "error must not be nil")
	require.Containsf(t, err.Error(), "rows error: ", "error: %s", err)
	require.Nil(t, names, "names must be nil")
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dbService := myDb.New(mockDB)

	require.NotNil(t, dbService, "dbService must not be nil")
}

func mockDBRows(names []string) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"name"})
	for _, name := range names {
		rows = rows.AddRow(name)
	}

	return rows
}

func TestDBGetUniqueNames(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dbService := myDb.DBService{DB: mockDB}

	for i, row := range testTable {
		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnRows(mockDBRowsUnique(row.names)).WillReturnError(row.errExpected)

		names, err := dbService.GetUniqueNames()

		if row.errExpected != nil {
			require.ErrorIs(t, err, row.errExpected, "row: %d, expected error: %w, actual error: %w", i,
				row.errExpected, err)
			require.Nil(t, names, "row: %d, names must be nil", i)

			continue
		}

		require.NoError(t, err, "row: %d, error must be nil", i)
		require.Equal(t, row.names, names, "row: %d, expected names: %s, actual names: %s",
			i, row.names, names)
	}
}

func TestGetUniqueNames_ScanError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dbService := myDb.DBService{DB: mockDB}

	// setting up the mock so that the query returns two columns
	rows := sqlmock.NewRows([]string{"name", "ages"}).AddRow("Ivan", "18")
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	names, err := dbService.GetUniqueNames()
	require.Error(t, err, "error must not be nil")
	require.Containsf(t, err.Error(), "rows scanning: ", "error: %s", err)
	require.Nil(t, names, "names must be nil")
}

func TestGetUniqueNames_RowError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dbService := myDb.DBService{DB: mockDB}

	// simulate a connection error when reading a row
	rows := sqlmock.NewRows([]string{"name"}).AddRow("Ivan").
		AddRow("Gena228").RowError(1, errNetwork)
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	names, err := dbService.GetUniqueNames()

	require.Error(t, err, "error must not be nil")
	require.Containsf(t, err.Error(), "rows error: ", "error: %s", err)
	require.Nil(t, names, "names must be nil")
}

func mockDBRowsUnique(names []string) *sqlmock.Rows {
	seen := make(map[string]bool)
	unique := []string{}

	for _, name := range names {
		if !seen[name] {
			seen[name] = true

			unique = append(unique, name)
		}
	}

	rows := sqlmock.NewRows([]string{"name"})
	for _, name := range unique {
		rows = rows.AddRow(name)
	}

	return rows
}
