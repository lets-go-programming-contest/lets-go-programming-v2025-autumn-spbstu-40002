package db

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

var (
	errQueryFailed = errors.New("query failed")
	errRowsError   = errors.New("rows error")
)

func TestNew(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := New(db)
	require.NotNil(t, store)
	require.Equal(t, db, store.DB)
}

func TestStore_GetNames(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*sqlmock.Sqlmock)
		want    []string
		wantErr bool
	}{
		{
			name: "success",
			setup: func(m *sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob")
				(*m).ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).WillReturnRows(rows)
			},
			want: []string{"Alice", "Bob"},
		},
		{
			name: "empty",
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
			name: "rows error",
			setup: func(m *sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					RowError(0, errRowsError)
				(*m).ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.setup(&mock)

			store := New(db)
			got, err := store.GetNames()

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

func TestStore_GetUniqueNames(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*sqlmock.Sqlmock)
		want    []string
		wantErr bool
	}{
		{
			name: "success",
			setup: func(m *sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob")
				(*m).ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).WillReturnRows(rows)
			},
			want: []string{"Alice", "Bob"},
		},
		{
			name: "query error",
			setup: func(m *sqlmock.Sqlmock) {
				(*m).ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).WillReturnError(errQueryFailed)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.setup(&mock)

			store := New(db)
			got, err := store.GetUniqueNames()

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
