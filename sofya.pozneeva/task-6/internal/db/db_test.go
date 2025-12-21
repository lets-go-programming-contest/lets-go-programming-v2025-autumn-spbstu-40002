package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	db "task-6/internal/db"
)

var (
	errDBDown = errors.New("db down")
	errRow    = errors.New("row error")
)

func TestDBService_GetNames_OK(t *testing.T) {
	t.Parallel()

	conn, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("alice").
		AddRow("bob")

	mock.ExpectQuery("^SELECT name FROM users$").WillReturnRows(rows)

	service := db.New(conn)

	got, err := service.GetNames()
	require.NoError(t, err)
	require.Equal(t, []string{"alice", "bob"}, got)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetNames_QueryError(t *testing.T) {
	t.Parallel()

	conn, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	mock.ExpectQuery("^SELECT name FROM users$").WillReturnError(errDBDown)

	service := db.New(conn)

	got, err := service.GetNames()
	require.Nil(t, got)
	require.Error(t, err)
	require.ErrorContains(t, err, "db query:")
	require.ErrorIs(t, err, errDBDown)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetNames_ScanError(t *testing.T) {
	t.Parallel()

	conn, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
	mock.ExpectQuery("^SELECT name FROM users$").WillReturnRows(rows)

	service := db.New(conn)

	got, err := service.GetNames()
	require.Nil(t, got)
	require.Error(t, err)
	require.ErrorContains(t, err, "rows scanning:")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetNames_RowsError(t *testing.T) {
	t.Parallel()

	conn, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("alice").
		RowError(0, errRow)

	mock.ExpectQuery("^SELECT name FROM users$").WillReturnRows(rows)

	service := db.New(conn)

	got, err := service.GetNames()
	require.Nil(t, got)
	require.Error(t, err)
	require.ErrorContains(t, err, "rows error:")
	require.ErrorIs(t, err, errRow)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetUniqueNames_OK(t *testing.T) {
	t.Parallel()

	conn, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("alice").
		AddRow("bob")

	mock.ExpectQuery("^SELECT DISTINCT name FROM users$").WillReturnRows(rows)

	service := db.New(conn)

	got, err := service.GetUniqueNames()
	require.NoError(t, err)
	require.Equal(t, []string{"alice", "bob"}, got)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetUniqueNames_QueryError(t *testing.T) {
	t.Parallel()

	conn, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	mock.ExpectQuery("^SELECT DISTINCT name FROM users$").WillReturnError(errDBDown)

	service := db.New(conn)

	got, err := service.GetUniqueNames()
	require.Nil(t, got)
	require.Error(t, err)
	require.ErrorContains(t, err, "db query:")
	require.ErrorIs(t, err, errDBDown)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetUniqueNames_ScanError(t *testing.T) {
	t.Parallel()

	conn, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
	mock.ExpectQuery("^SELECT DISTINCT name FROM users$").WillReturnRows(rows)

	service := db.New(conn)

	got, err := service.GetUniqueNames()
	require.Nil(t, got)
	require.Error(t, err)
	require.ErrorContains(t, err, "rows scanning:")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetUniqueNames_RowsError(t *testing.T) {
	t.Parallel()

	conn, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("alice").
		RowError(0, errRow)

	mock.ExpectQuery("^SELECT DISTINCT name FROM users$").WillReturnRows(rows)

	service := db.New(conn)

	got, err := service.GetUniqueNames()
	require.Nil(t, got)
	require.Error(t, err)
	require.ErrorContains(t, err, "rows error:")
	require.ErrorIs(t, err, errRow)

	require.NoError(t, mock.ExpectationsWereMet())
}
