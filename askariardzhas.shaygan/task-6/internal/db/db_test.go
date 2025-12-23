package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/XShaygaND/task-6/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	errQuery = errors.New("query execution error")
	errRows  = errors.New("row processing error")
)

func TestNew(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	mock.ExpectClose()

	s := db.New(mockDB)

	assert.Equal(t, mockDB, s.DB)

	err = mockDB.Close()
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		setup    func(mock sqlmock.Sqlmock)
		expected []string
		errMsg   string
	}{
		{
			name: "success multiple rows",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob")
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expected: []string{"Alice", "Bob"},
		},
		{
			name: "success empty",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expected: nil,
		},
		{
			name: "query error",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT name FROM users").WillReturnError(errQuery)
			},
			errMsg: "query execution failed: query execution error",
		},
		{
			name: "scan error",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			errMsg: "scanning row failed",
		},
		{
			name: "rows error after iteration",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob").
					RowError(1, errRows)
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			errMsg: "rows iteration error: row processing error",
		},
		{
			name: "rows error no rows",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("dummy").
					RowError(0, errRows)
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			errMsg: "rows iteration error: row processing error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)

			tc.setup(mock)
			mock.ExpectClose()

			s := db.New(mockDB)
			got, err := s.GetNames()

			closeErr := mockDB.Close()
			require.NoError(t, closeErr)

			if tc.errMsg != "" {
				require.Error(t, err)

				assert.Contains(t, err.Error(), tc.errMsg)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)

				assert.Equal(t, tc.expected, got)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		setup    func(mock sqlmock.Sqlmock)
		expected []string
		errMsg   string
	}{
		{
			name: "success multiple rows",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expected: []string{"Alice", "Bob"},
		},
		{
			name: "success empty",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expected: nil,
		},
		{
			name: "query error",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errQuery)
			},
			errMsg: "query execution failed: query execution error",
		},
		{
			name: "scan error",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			errMsg: "scanning row failed",
		},
		{
			name: "rows error after iteration",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob").
					RowError(1, errRows)
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			errMsg: "rows iteration error: row processing error",
		},
		{
			name: "rows error no rows",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("dummy").
					RowError(0, errRows)
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			errMsg: "rows iteration error: row processing error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)

			tc.setup(mock)
			mock.ExpectClose()

			s := db.New(mockDB)
			got, err := s.GetUniqueNames()

			closeErr := mockDB.Close()
			require.NoError(t, closeErr)

			if tc.errMsg != "" {
				require.Error(t, err)

				assert.Contains(t, err.Error(), tc.errMsg)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)

				assert.Equal(t, tc.expected, got)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
