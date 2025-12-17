package db_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Svyatoslav2324/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

func newMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, db.DBService) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)

	service := db.New(dbMock)
	return dbMock, mock, service
}

func TestGetNames_Success(t *testing.T) {
	dbMock, mock, service := newMockDB(t)
	defer dbMock.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob")

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnRows(rows)

	names, err := service.GetNames()

	require.NoError(t, err)
	require.Equal(t, []string{"Alice", "Bob"}, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_QueryError(t *testing.T) {
	dbMock, mock, service := newMockDB(t)
	defer dbMock.Close()

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnError(errors.New("query failed"))

	names, err := service.GetNames()

	require.Error(t, err)
	require.Nil(t, names)
	require.Contains(t, err.Error(), "db query")
}

func TestGetNames_ScanError(t *testing.T) {
	dbMock, mock, service := newMockDB(t)
	defer dbMock.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow(123) // неверный тип → Scan error

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnRows(rows)

	names, err := service.GetNames()

	require.Error(t, err)
	require.Nil(t, names)
	require.Contains(t, err.Error(), "rows scanning")
}

func TestGetNames_RowsError(t *testing.T) {
	dbMock, mock, service := newMockDB(t)
	defer dbMock.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errors.New("rows error"))

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnRows(rows)

	names, err := service.GetNames()

	require.Error(t, err)
	require.Nil(t, names)
	require.Contains(t, err.Error(), "rows error")
}

func TestGetUniqueNames_Success(t *testing.T) {
	dbMock, mock, service := newMockDB(t)
	defer dbMock.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob")

	mock.ExpectQuery("SELECT DISTINCT name FROM users").
		WillReturnRows(rows)

	names, err := service.GetUniqueNames()

	require.NoError(t, err)
	require.Equal(t, []string{"Alice", "Bob"}, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_QueryError(t *testing.T) {
	dbMock, mock, service := newMockDB(t)
	defer dbMock.Close()

	mock.ExpectQuery("SELECT DISTINCT name FROM users").
		WillReturnError(errors.New("query failed"))

	names, err := service.GetUniqueNames()

	require.Error(t, err)
	require.Nil(t, names)
	require.Contains(t, err.Error(), "db query")
}

func TestGetUniqueNames_ScanError(t *testing.T) {
	dbMock, mock, service := newMockDB(t)
	defer dbMock.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow(123)

	mock.ExpectQuery("SELECT DISTINCT name FROM users").
		WillReturnRows(rows)

	names, err := service.GetUniqueNames()

	require.Error(t, err)
	require.Nil(t, names)
	require.Contains(t, err.Error(), "rows scanning")
}

func TestGetUniqueNames_RowsError(t *testing.T) {
	dbMock, mock, service := newMockDB(t)
	defer dbMock.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errors.New("rows error"))

	mock.ExpectQuery("SELECT DISTINCT name FROM users").
		WillReturnRows(rows)

	names, err := service.GetUniqueNames()

	require.Error(t, err)
	require.Nil(t, names)
	require.Contains(t, err.Error(), "rows error")
}
