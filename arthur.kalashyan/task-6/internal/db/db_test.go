package db_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Expeline/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

var (
	errQueryFailed = errors.New("query failed")
	errRowsError   = errors.New("rows error")
	errCloseError  = errors.New("close error")
	errFail        = errors.New("fail")
)

type dbTestEnv struct {
	Mock    sqlmock.Sqlmock
	Service db.DBService
}

func setupDB(t *testing.T) dbTestEnv {
	t.Helper()

	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, dbConn.Close())
	})

	return dbTestEnv{
		Mock:    mock,
		Service: db.New(dbConn),
	}
}

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(sqlmock.Sqlmock)
		want    []string
		wantErr bool
	}{
		{
			name: "ok_multiple_rows",
			setup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("a").
					AddRow("b")
				m.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			want: []string{"a", "b"},
		},
		{
			name: "empty_result",
			setup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				m.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			want: []string{},
		},
		{
			name: "query_error",
			setup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT name FROM users").
					WillReturnError(errQueryFailed)
			},
			wantErr: true,
		},
		{
			name: "scan_error",
			setup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(123)
				m.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name: "rows_error_after_iteration",
			setup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("a").
					RowError(0, errRowsError)
				m.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name: "rows_close_error",
			setup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("a").
					CloseError(errCloseError)
				m.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			env := setupDB(t)
			tt.setup(env.Mock)

			res, err := env.Service.GetNames()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, res)
			}

			require.NoError(t, env.Mock.ExpectationsWereMet())
		})
	}
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(sqlmock.Sqlmock)
		want    []string
		wantErr bool
	}{
		{
			name: "ok",
			setup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("a").
					AddRow("b")
				m.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			want: []string{"a", "b"},
		},
		{
			name: "empty_result",
			setup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				m.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			want: []string{},
		},
		{
			name: "query_error",
			setup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnError(errFail)
			},
			wantErr: true,
		},
		{
			name: "scan error",
			setup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				m.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name: "rows_error",
			setup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("a").
					RowError(0, errRowsError)
				m.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name: "rows_close_error",
			setup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("a").
					CloseError(errCloseError)
				m.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			env := setupDB(t)
			tt.setup(env.Mock)

			res, err := env.Service.GetUniqueNames()
			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, res)
			require.NoError(t, env.Mock.ExpectationsWereMet())
		})
	}
}
