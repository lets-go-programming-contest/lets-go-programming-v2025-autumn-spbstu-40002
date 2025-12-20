package db

import (
	"database/sql"
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

type sqlDBWrapper struct {
	db *sql.DB
}

func (w sqlDBWrapper) Query(query string, args ...any) (Rows, error) {
	return w.db.Query(query, args...)
}

type fakeRows struct {
	names    []string
	idx      int
	scanErr  bool
	errValue error
	closeErr error
}

func (r *fakeRows) Next() bool {
	return r.idx < len(r.names)
}

func (r *fakeRows) Scan(dest ...any) error {
	if r.scanErr {
		return errors.New("scan failed")
	}
	if len(dest) == 0 {
		return errors.New("no dest")
	}
	ptr, ok := dest[0].(*string)
	if !ok {
		return errors.New("bad dest type")
	}
	*ptr = r.names[r.idx]
	r.idx++
	return nil
}

func (r *fakeRows) Err() error {
	return r.errValue
}

func (r *fakeRows) Close() error {
	return r.closeErr
}

type fakeDB struct {
	rows     Rows
	queryErr error
}

func (f fakeDB) Query(query string, args ...any) (Rows, error) {
	if f.queryErr != nil {
		return nil, f.queryErr
	}
	return f.rows, nil
}

func TestNew(t *testing.T) {
	dbSQL, _, err := sqlmock.New()
	require.NoError(t, err)
	defer dbSQL.Close()

	w := sqlDBWrapper{db: dbSQL}
	store := New(w)
	require.NotNil(t, store)
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
			dbSQL, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer dbSQL.Close()

			tt.setup(&mock)

			store := New(sqlDBWrapper{db: dbSQL})
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

func TestStore_GetNames_CloseError(t *testing.T) {
	fr := &fakeRows{
		names:    []string{},
		closeErr: errors.New("close failed"),
	}
	store := New(fakeDB{rows: fr})
	_, err := store.GetNames()
	require.Error(t, err)
	require.Regexp(t, regexp.MustCompile(`close: .*close failed`), err.Error())
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbSQL, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer dbSQL.Close()

			tt.setup(&mock)

			store := New(sqlDBWrapper{db: dbSQL})
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

func TestStore_GetUniqueNames_CloseError(t *testing.T) {
	fr := &fakeRows{
		names:    []string{},
		closeErr: errors.New("close failed"),
	}
	store := New(fakeDB{rows: fr})
	_, err := store.GetUniqueNames()
	require.Error(t, err)
	require.Regexp(t, regexp.MustCompile(`close: .*close failed`), err.Error())
}

func TestStore_GetUniqueNames_ScanError(t *testing.T) {
	fr := &fakeRows{
		names:   []string{"Alice"},
		scanErr: true,
	}

	store := New(fakeDB{rows: fr})

	_, err := store.GetUniqueNames()
	require.Error(t, err)
	require.Contains(t, err.Error(), "scan:")
}

func TestStore_GetUniqueNames_RowsError(t *testing.T) {
	fr := &fakeRows{
		names:    []string{"Alice"},
		errValue: errRowsError,
	}

	store := New(fakeDB{rows: fr})

	_, err := store.GetUniqueNames()
	require.Error(t, err)
	require.Contains(t, err.Error(), "rows:")
}

func TestStore_GetUniqueNames_CloseDoesNotOverrideScanError(t *testing.T) {
	fr := &fakeRows{
		names:    []string{"Alice"},
		scanErr:  true,
		closeErr: errors.New("close failed"),
	}

	store := New(fakeDB{rows: fr})

	_, err := store.GetUniqueNames()
	require.Error(t, err)
	require.Contains(t, err.Error(), "scan:")
}

func TestStore_GetUniqueNames_CloseDoesNotOverrideRowsError(t *testing.T) {
	fr := &fakeRows{
		names:    []string{"Alice"},
		errValue: errRowsError,
		closeErr: errors.New("close failed"),
	}

	store := New(fakeDB{rows: fr})

	_, err := store.GetUniqueNames()
	require.Error(t, err)
	require.Contains(t, err.Error(), "rows:")
}
