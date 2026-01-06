package db_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/hehemka/task-6/internal/db"
)

var (
	errQuery = errors.New("query error")
	errRow   = errors.New("row error")
)

func TestDBServiceNew(t *testing.T) {
	t.Parallel()

	conn, _, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	svc := db.New(conn)
	require.NotNil(t, svc)
	require.Same(t, conn, svc.DB)
}

func TestDBServiceGetNames(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name   string
		setup  func(sqlmock.Sqlmock)
		assert func(*testing.T, []string, error)
	}

	tests := []testCase{
		{
			name: "success",
			setup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta(
					"SELECT name FROM users",
				)).WillReturnRows(
					sqlmock.NewRows([]string{"name"}).
						AddRow("Alice").
						AddRow("Bob"),
				)
			},
			assert: func(t *testing.T, result []string, err error) {
				require.NoError(t, err)
				require.Equal(t, []string{"Alice", "Bob"}, result)
			},
		},
		{
			name: "empty",
			setup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta(
					"SELECT name FROM users",
				)).WillReturnRows(
					sqlmock.NewRows([]string{"name"}),
				)
			},
			assert: func(t *testing.T, result []string, err error) {
				require.NoError(t, err)
				require.Empty(t, result)
			},
		},
		{
			name: "query error",
			setup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta(
					"SELECT name FROM users",
				)).WillReturnError(errQuery)
			},
			assert: func(t *testing.T, result []string, err error) {
				require.Error(t, err)
				require.Nil(t, result)
			},
		},
		{
			name: "scan error",
			setup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta(
					"SELECT name FROM users",
				)).WillReturnRows(
					sqlmock.NewRows([]string{"name"}).AddRow(nil),
				)
			},
			assert: func(t *testing.T, result []string, err error) {
				require.Error(t, err)
				require.Nil(t, result)
			},
		},
		{
			name: "rows error",
			setup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta(
					"SELECT name FROM users",
				)).WillReturnRows(
					sqlmock.NewRows([]string{"name"}).
						AddRow("Alice").
						RowError(0, errRow),
				)
			},
			assert: func(t *testing.T, result []string, err error) {
				require.Error(t, err)
				require.Nil(t, result)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			conn, mock, err := sqlmock.New()
			require.NoError(t, err)
			t.Cleanup(func() { _ = conn.Close() })

			svc := db.New(conn)
			tt.setup(mock)

			result, err := svc.GetNames()
			tt.assert(t, result, err)

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBServiceGetUniqueNames(t *testing.T) {
	t.Parallel()

	conn, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	svc := db.New(conn)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT DISTINCT name FROM users",
	)).WillReturnRows(
		sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob"),
	)

	names, err := svc.GetUniqueNames()

	require.NoError(t, err)
	require.ElementsMatch(t, []string{"Alice", "Bob"}, names)
	require.NoError(t, mock.ExpectationsWereMet())
}
