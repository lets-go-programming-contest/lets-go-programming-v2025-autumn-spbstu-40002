package db

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestGetNames_SuccessMultipleRows(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		if cerr := mockDB.Close(); cerr != nil {
			t.Logf("mockDB.Close: %v", cerr)
		}
	}()
	rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").AddRow("Bob").AddRow("Carol")
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
	s := New(mockDB)
	names, err := s.GetNames()
	require.NoError(t, err)
	require.Equal(t, []string{"Alice", "Bob", "Carol"}, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_EmptyRows(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		if cerr := mockDB.Close(); cerr != nil {
			t.Logf("mockDB.Close: %v", cerr)
		}
	}()
	rows := sqlmock.NewRows([]string{"name"})
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
	s := New(mockDB)
	names, err := s.GetNames()
	require.NoError(t, err)
	require.Len(t, names, 0)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_QueryError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		if cerr := mockDB.Close(); cerr != nil {
			t.Logf("mockDB.Close: %v", cerr)
		}
	}()
	mock.ExpectQuery("SELECT name FROM users").WillReturnError(sql.ErrConnDone)
	s := New(mockDB)
	_, err = s.GetNames()
	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_ScanError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		if cerr := mockDB.Close(); cerr != nil {
			t.Logf("mockDB.Close: %v", cerr)
		}
	}()
	rows := sqlmock.NewRows([]string{"name"}).AddRow("ok").AddRow("bad")
	rows.RowError(1, errors.New("scan fail"))
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
	s := New(mockDB)
	_, err = s.GetNames()
	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_RowsErr(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		if cerr := mockDB.Close(); cerr != nil {
			t.Logf("mockDB.Close: %v", cerr)
		}
	}()
	rows := sqlmock.NewRows([]string{"name"}).AddRow("a")
	rows.RowError(0, errors.New("row iteration error"))
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
	s := New(mockDB)
	_, err = s.GetNames()
	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_CloseErrorBranch(t *testing.T) {
	old := closeRows
	defer func() { closeRows = old }()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		if cerr := mockDB.Close(); cerr != nil {
			t.Log(cerr)
		}
	}()

	closeRows = func(r *sql.Rows) error {
		return errors.New("close error")
	}

	rows := sqlmock.NewRows([]string{"name"}).AddRow("A")
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	s := New(mockDB)
	_, err = s.GetNames()
	require.Error(t, err)
}
