package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Svyatoslav2324/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

var (
	ErrRowsError  = errors.New("rows error")
	ErrQueryError = errors.New("query error")
)

func mockDBRows(names []string) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"name"})
	for _, name := range names {
		rows = rows.AddRow(name)
	}

	return rows
}

type getNamesFunc func() ([]string, error)

type getNamesTestCase struct {
	name          string
	rows          *sqlmock.Rows
	expected      []string
	expectedError bool
	queryError    error
}

func getNamesTestCases() []getNamesTestCase {
	return []getNamesTestCase{
		{
			name:     "success",
			rows:     mockDBRows([]string{"Bob", "Steve", "Svetozar"}),
			expected: []string{"Bob", "Steve", "Svetozar"},
		},
		{
			name: "rows error",
			rows: sqlmock.NewRows([]string{"name"}).
				AddRow("Bob").
				RowError(0, ErrRowsError),
			expectedError: true,
		},
		{
			name:          "scan error",
			rows:          sqlmock.NewRows([]string{"name"}).AddRow(nil),
			expectedError: true,
		},
		{
			name:     "no rows",
			rows:     sqlmock.NewRows([]string{"name"}),
			expected: nil,
		},
		{
			name:          "query error",
			queryError:    ErrQueryError,
			expectedError: true,
		},
	}
}

func testGetNamesMethod(
	t *testing.T,
	query string,
	method func(db.DBService) getNamesFunc,
) {
	t.Helper()

	for _, testCase := range getNamesTestCases() {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)

			service := db.DBService{DB: mockDB}

			expect := mock.ExpectQuery(query)

			if testCase.queryError != nil {
				expect.WillReturnError(testCase.queryError)
			} else {
				expect.WillReturnRows(testCase.rows)
			}

			names, err := method(service)()

			if testCase.expectedError {
				require.Error(t, err)
				require.Nil(t, names)

				return
			}

			require.NoError(t, err)
			require.Equal(t, testCase.expected, names)
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	/*
		Насколько я понял, этот тест не имеет смысла.
		Конструктор не имеет вообще никого ветвления -> нечему ломаться.
		Но тест есть, чтобы обеспечить покрытие 100%.
	*/
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dbService := db.New(mockDB)

	require.Equal(t, mockDB, dbService.DB,
		"expected: %s, got: %s", mockDB, dbService.DB)
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	testGetNamesMethod(
		t,
		"SELECT name FROM users",
		func(service db.DBService) getNamesFunc {
			return service.GetNames
		})
}

func TestGetUniqueNames(t *testing.T) {
	t.Parallel()

	testGetNamesMethod(
		t,
		"SELECT DISTINCT name FROM users",
		func(service db.DBService) getNamesFunc {
			return service.GetUniqueNames
		})
}
