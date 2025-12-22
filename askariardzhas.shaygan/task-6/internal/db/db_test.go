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
	errRow   = errors.New("row processing error")
)

func TestCreate(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	mock.ExpectClose()

	client := db.Create(mockDB)

	assert.Equal(t, mockDB, client.Conn)

	err = mockDB.Close()
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBClient_GetAllNames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		mockFunc func(mock sqlmock.Sqlmock)
		expected []string
		errMsg   string
	}{
		{
			name: "multiple results",
			mockFunc: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob")
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expected: []string{"Alice", "Bob"},
		},
		{
			name: "no results",
			mockFunc: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expected: nil,
		},
		{
			name: "query failure",
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT name FROM users").WillReturnError(errQuery)
			},
			errMsg: "query execution failed: query execution error",
		},
		{
			name: "scan failure",
			mockFunc: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			errMsg: "scanning row failed",
		},
		{
			name: "row error during iteration",
			mockFunc: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob").
					RowError(1, errRow)
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			errMsg: "rows iteration error: row processing error",
		},
		{
			name: "row error with no rows",
			mockFunc: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("test").
					RowError(0, errRow)
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

			tc.mockFunc(mock)
			mock.ExpectClose()

			client := db.Create(mockDB)
			result, err := client.GetAllNames()

			closeErr := mockDB.Close()
			require.NoError(t, closeErr)

			if tc.errMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errMsg)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBClient_GetDistinctNames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		mockFunc func(mock sqlmock.Sqlmock)
		expected []string
		errMsg   string
	}{
		{
			name: "multiple distinct results",
			mockFunc: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expected: []string{"Alice", "Bob"},
		},
		{
			name: "no distinct results",
			mockFunc: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expected: nil,
		},
		{
			name: "distinct query failure",
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errQuery)
			},
			errMsg: "query execution failed: query execution error",
		},
		{
			name: "distinct scan failure",
			mockFunc: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			errMsg: "scanning row failed",
		},
		{
			name: "distinct row error during iteration",
			mockFunc: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob").
					RowError(1, errRow)
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			errMsg: "rows iteration error: row processing error",
		},
		{
			name: "distinct row error with no rows",
			mockFunc: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("sample").
					RowError(0, errRow)
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

			tc.mockFunc(mock)
			mock.ExpectClose()

			client := db.Create(mockDB)
			result, err := client.GetDistinctNames()

			closeErr := mockDB.Close()
			require.NoError(t, closeErr)

			if tc.errMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errMsg)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
