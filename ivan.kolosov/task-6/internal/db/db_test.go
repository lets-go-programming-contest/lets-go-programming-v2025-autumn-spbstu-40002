package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/InsomniaDemon/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

var errExpected = errors.New("errExpected")

type rowTestDB struct {
	mockRows      *sqlmock.Rows
	expectedError error
	expectations  []string
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)

	require.Equal(t, mockDB, service.DB)
}

func TestGetNames(t *testing.T) {
	testTable := []rowTestDB{
		{
			mockRows:     sqlmock.NewRows([]string{"name"}).AddRow("Dasha").AddRow("Daria").AddRow("Dashka"),
			expectations: []string{"Dasha", "Daria", "Dashka"},
		},
		{
			mockRows:      sqlmock.NewRows([]string{"name"}).AddRow("Dasha").AddRow("Daria").AddRow("Dashka"),
			expectedError: errExpected,
		},
	}

	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	defer mockDB.Close()

	dbService := db.DBService{DB: mockDB}

	for i, testCase := range testTable {
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(testCase.mockRows).
			WillReturnError(testCase.expectedError)

		names, err := dbService.GetNames()

		if testCase.expectedError != nil {
			require.ErrorIs(t, err, testCase.expectedError, "row: %d, expected error: %w, actual error: %w", i,
				testCase.expectedError, err)
			require.Nil(t, names, "row: %d, names must be nil", i)

			continue
		}

		require.NoError(t, err, "row: %d, error must be nil", i)
		require.Equal(t, testCase.expectations, names, "row: %d, expected names: %s, actual names: %s", i,
			testCase.expectations, names)
	}

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))

	names, err := dbService.GetNames()

	require.ErrorContains(t, err, "rows scanning")
	require.Nil(t, names, "row: %d, names must be nil")

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).
		AddRow("DASHKENT").AddRow("Dashka").RowError(1, errExpected))

	names, err = dbService.GetNames()
	require.ErrorContains(t, err, "rows error")
	require.Nil(t, names, "row: %d, names must be nil")
}

func TestGetUniqueNames(t *testing.T) {
	testTable := []rowTestDB{
		{
			mockRows:     sqlmock.NewRows([]string{"name"}).AddRow("Dasha").AddRow("Daria").AddRow("Dashka"),
			expectations: []string{"Dasha", "Daria", "Dashka"},
		},
		{
			mockRows:      sqlmock.NewRows([]string{"name"}).AddRow("Dasha").AddRow("Daria").AddRow("Dashka"),
			expectedError: errExpected,
		},
	}

	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	defer mockDB.Close()

	dbService := db.DBService{DB: mockDB}

	for i, testCase := range testTable {
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(testCase.mockRows).
			WillReturnError(testCase.expectedError)

		names, err := dbService.GetUniqueNames()

		if testCase.expectedError != nil {
			require.ErrorIs(t, err, testCase.expectedError, "row: %d, expected error: %w, actual error: %w", i,
				testCase.expectedError, err)
			require.Nil(t, names, "row: %d, names must be nil", i)

			continue
		}

		require.NoError(t, err, "row: %d, error must be nil", i)
		require.Equal(t, testCase.expectations, names, "row: %d, expected names: %s, actual names: %s", i,
			testCase.expectations, names)
	}

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))

	names, err := dbService.GetUniqueNames()

	require.ErrorContains(t, err, "rows scanning")
	require.Nil(t, names, "row: %d, names must be nil")

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).
		AddRow("DASHKENT").AddRow("Dashka").RowError(1, errExpected))

	names, err = dbService.GetUniqueNames()
	require.ErrorContains(t, err, "rows error")
	require.Nil(t, names, "row: %d, names must be nil")
}
