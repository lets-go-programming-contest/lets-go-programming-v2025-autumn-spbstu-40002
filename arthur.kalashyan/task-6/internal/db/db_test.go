package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

type dbTestEnv struct {
	Mock    sqlmock.Sqlmock
	Service DBService
}

func setupDB(t *testing.T) dbTestEnv {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})

	return dbTestEnv{
		Mock:    mock,
		Service: New(db),
	}
}

func TestDBService_GetNames(t *testing.T) {
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
				m.ExpectQuery("SELECT name FROM users").
					WillReturnRows(rows)
			},
			want: []string{"a", "b"},
		},
		{
			name: "empty_result",
			setup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				m.ExpectQuery("SELECT name FROM users").
					WillReturnRows(rows)
			},
			want: []string{},
		},
		{
			name: "query_error",
			setup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT name FROM users").
					WillReturnError(errors.New("query failed"))
			},
			wantErr: true,
		},
		{
			name: "scan_error",
			setup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow(123)
				m.ExpectQuery("SELECT name FROM users").
					WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name: "rows_error_after_iteration",
			setup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("a").
					RowError(0, errors.New("rows error"))
				m.ExpectQuery("SELECT name FROM users").
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
				m.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnRows(rows)
			},
			want: []string{"a", "b"},
		},
		{
			name: "empty_result",
			setup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				m.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnRows(rows)
			},
			want: []string{},
		},
		{
			name: "query_error",
			setup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnError(errors.New("fail"))
			},
			wantErr: true,
		},
		{
			name: "rows_error",
			setup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("a").
					RowError(0, errors.New("rows error"))
				m.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
