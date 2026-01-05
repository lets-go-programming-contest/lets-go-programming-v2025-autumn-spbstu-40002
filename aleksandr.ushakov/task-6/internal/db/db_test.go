package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

type rowTestDb struct {
	names       []string
	errExpected error
}

var testTable = []rowTestDb{
	{
		names: []string{"Ivan", "Gena228"},
	},
	{
		names:       nil,
		errExpected: errors.New("ExpectedError"),
	},
}

func TestGetNames(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbService := DBService{DB: mockDb}

	for i, row := range testTable {
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(mockDbRows(row.names)).
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
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbService := DBService{DB: mockDb}
	//setting up the mock so that the query returns two columns
	rows := sqlmock.NewRows([]string{"name", "ages"}).AddRow("Ivan", "18")
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	names, err := dbService.GetNames()
	require.Error(t, err, "error must not be nil")
	require.Containsf(t, err.Error(), "rows scanning: ", "error: %s", err)
	require.Nil(t, names, "names must be nil")
}

func TestGetNames_RowError(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbService := DBService{DB: mockDb}

	// simulate a connection error when reading a row
	rows := sqlmock.NewRows([]string{"name"}).AddRow("Ivan").
		AddRow("Gena228").RowError(1, errors.New("network problem"))
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
	names, err := dbService.GetNames()

	require.Error(t, err, "error must not be nil")
	require.Containsf(t, err.Error(), "rows error: ", "error: %s", err)
	require.Nil(t, names, "names must be nil")

}

func TestNew(t *testing.T) {
	mockDb, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbService := New(mockDb)
	require.NotNil(t, dbService, "dbService must not be nil")

}

func mockDbRows(names []string) *sqlmock.Rows {

	rows := sqlmock.NewRows([]string{"name"})
	for _, name := range names {
		rows = rows.AddRow(name)
	}
	return rows
}

func TestDBGetUniqueNamesNames(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbService := DBService{DB: mockDb}
	for i, row := range testTable {
		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnRows(mockDbRowsUnique(row.names)).WillReturnError(row.errExpected)

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
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbService := DBService{DB: mockDb}

	//setting up the mock so that the query returns two columns
	rows := sqlmock.NewRows([]string{"name", "ages"}).AddRow("Ivan", "18")
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	names, err := dbService.GetUniqueNames()
	require.Error(t, err, "error must not be nil")
	require.Containsf(t, err.Error(), "rows scanning: ", "error: %s", err)
	require.Nil(t, names, "names must be nil")
}

func TestGetUniqueNames_RowError(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbService := DBService{DB: mockDb}

	// simulate a connection error when reading a row
	rows := sqlmock.NewRows([]string{"name"}).AddRow("Ivan").
		AddRow("Gena228").RowError(1, errors.New("network problem"))
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
	names, err := dbService.GetUniqueNames()

	require.Error(t, err, "error must not be nil")
	require.Containsf(t, err.Error(), "rows error: ", "error: %s", err)
	require.Nil(t, names, "names must be nil")

}

func mockDbRowsUnique(names []string) *sqlmock.Rows {
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
