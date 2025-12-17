package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Exam-Play/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

const (
	queryGetNames       = "SELECT name FROM users"
	queryGetUniqueNames = "SELECT DISTINCT name FROM users"
)

var errExpected = errors.New("errExpected")

type testCase struct {
	name        string
	mockRows    *sqlmock.Rows
	expectError bool
	names       []string
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
	t.Parallel()

	tests := []testCase{
		{
			name:        "success - multiple names",
			mockRows:    sqlmock.NewRows([]string{"name"}).AddRow("Ivan").AddRow("Gena228"),
			expectError: false,
			names:       []string{"Ivan", "Gena228"},
		},
		{
			name:        "error - query failed",
			mockRows:    nil,
			expectError: true,
			names:       nil,
		},
	}

	t.Run("rows error after first row", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Ivan").
			AddRow("Petr").
			RowError(1, errExpected)

		mock.ExpectQuery(queryGetNames).
			WillReturnRows(rows)

		dbService := db.New(mockDB)
		names, err := dbService.GetNames()

		require.ErrorContains(t, err, "rows error")
		require.Nil(t, names)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			dbService := db.New(mockDB)

			if tt.expectError {
				mock.ExpectQuery(queryGetNames).WillReturnError(errExpected)
			} else {
				mock.ExpectQuery(queryGetNames).WillReturnRows(tt.mockRows)
			}

			names, err := dbService.GetNames()

			if tt.expectError {
				require.ErrorIs(t, err, errExpected)
				require.Nil(t, names)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.names, names)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}

	t.Run("scan error", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		mock.ExpectQuery(queryGetNames).
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))

		dbService := db.New(mockDB)
		names, err := dbService.GetNames()

		require.ErrorContains(t, err, "rows scanning")
		require.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Ivan").
			RowError(1, errExpected)

		mock.ExpectQuery(queryGetNames).WillReturnRows(rows)

		dbService := db.New(mockDB)
		names, err := dbService.GetNames()

		require.ErrorContains(t, err, "rows error")
		require.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetUniqueNames(t *testing.T) {
	t.Parallel()

	tests := []testCase{
		{
			name:        "success - multiple names",
			mockRows:    sqlmock.NewRows([]string{"name"}).AddRow("Ivan").AddRow("Gena228"),
			expectError: false,
			names:       []string{"Ivan", "Gena228"},
		},
		{
			name:        "error - query failed",
			mockRows:    nil,
			expectError: true,
			names:       nil,
		},
	}

	t.Run("rows error after first row", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Ivan").
			AddRow("Ivan").
			RowError(1, errExpected)

		mock.ExpectQuery(queryGetUniqueNames).
			WillReturnRows(rows)

		dbService := db.New(mockDB)
		names, err := dbService.GetUniqueNames()

		require.ErrorContains(t, err, "rows error")
		require.Nil(t, names)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			dbService := db.New(mockDB)

			if tt.expectError {
				mock.ExpectQuery(queryGetUniqueNames).WillReturnError(errExpected)
			} else {
				mock.ExpectQuery(queryGetUniqueNames).WillReturnRows(tt.mockRows)
			}

			names, err := dbService.GetUniqueNames()

			if tt.expectError {
				require.ErrorIs(t, err, errExpected)
				require.Nil(t, names)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.names, names)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}

	t.Run("scan error", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		mock.ExpectQuery(queryGetUniqueNames).
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))

		dbService := db.New(mockDB)
		names, err := dbService.GetUniqueNames()

		require.ErrorContains(t, err, "rows scanning")
		require.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Ivan").
			RowError(1, errExpected)

		mock.ExpectQuery(queryGetUniqueNames).WillReturnRows(rows)

		dbService := db.New(mockDB)
		names, err := dbService.GetUniqueNames()

		require.ErrorContains(t, err, "rows error")
		require.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
