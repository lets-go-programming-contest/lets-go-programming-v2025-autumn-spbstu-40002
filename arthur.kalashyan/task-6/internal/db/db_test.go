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
)

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(*sqlmock.Sqlmock)
		want    []string
		wantErr bool
	}{
		{
			name: "success with data",
			setup: func(m *sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").AddRow("Bob")
				(*m).ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).WillReturnRows(rows)
			},
			want: []string{"Alice", "Bob"},
		},
		{
			name: "success empty",
			setup: func(m *sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				(*m).ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).WillReturnRows(rows)
			},
			want: []string{},
		},
		{
			name: "query error",
			setup: func(m *sqlmock.Sqlmock) {
				(*m).ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).WillReturnError(errQueryFailed)
			},
			wantErr: true,
		},
		{
			name: "scan error",
			setup: func(m *sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				(*m).ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name: "rows error after iteration",
			setup: func(m *sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").RowError(0, errRowsError)
				(*m).ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name: "rows close error",
			setup: func(m *sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").CloseError(errCloseError)
				(*m).ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			conn, mock, err := sqlmock.New()
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, conn.Close()) })

			svc := db.New(conn)

			tt.setup(&mock)

			got, err := svc.GetNames()
			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(*sqlmock.Sqlmock)
		want    []string
		wantErr bool
	}{
		{
			name: "success with data",
			setup: func(m *sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").AddRow("Bob")
				(*m).ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).WillReturnRows(rows)
			},
			want: []string{"Alice", "Bob"},
		},
		{
			name: "success empty",
			setup: func(m *sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				(*m).ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).WillReturnRows(rows)
			},
			want: []string{},
		},
		{
			name: "query error",
			setup: func(m *sqlmock.Sqlmock) {
				(*m).ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).WillReturnError(errQueryFailed)
			},
			wantErr: true,
		},
		{
			name: "scan error",
			setup: func(m *sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				(*m).ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name: "rows error after iteration",
			setup: func(m *sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").RowError(0, errRowsError)
				(*m).ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name: "rows close error",
			setup: func(m *sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").CloseError(errCloseError)
				(*m).ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			conn, mock, err := sqlmock.New()
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, conn.Close()) })

			svc := db.New(conn)

			tt.setup(&mock)

			got, err := svc.GetUniqueNames()
			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
