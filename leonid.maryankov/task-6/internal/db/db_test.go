package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/leonid.maryankov/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)

	service := db.New(mockDB)
	require.NotNil(t, service)
}

func TestGetNames(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mock.ExpectationsWereMet()

	service := db.New(mockDB)

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob")

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnRows(rows)

	names, err := service.GetNames()

	require.NoError(t, err)
	require.Equal(t, []string{"Alice", "Bob"}, names)
}

func TestGetNames_EmptyResult(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mock.ExpectationsWereMet()

	service := db.New(mockDB)

	rows := sqlmock.NewRows([]string{"name"})

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnRows(rows)

	names, err := service.GetNames()

	require.NoError(t, err)
	require.Empty(t, names)
}

func TestGetNames_QueryError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mock.ExpectationsWereMet()

	service := db.New(mockDB)

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnError(errors.New("db error"))

	names, err := service.GetNames()

	require.Error(t, err)
	require.Nil(t, names)
}

func TestGetNames_ScanError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mock.ExpectationsWereMet()

	service := db.New(mockDB)

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errors.New("scan error"))

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnRows(rows)

	names, err := service.GetNames()

	require.Error(t, err)
	require.Nil(t, names)
}

func TestGetNames_RowsError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mock.ExpectationsWereMet()

	service := db.New(mockDB)

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		CloseError(errors.New("rows error"))

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnRows(rows)

	names, err := service.GetNames()

	require.Error(t, err)
	require.Nil(t, names)
}

func TestGetUniqueNames(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mock.ExpectationsWereMet()

	service := db.New(mockDB)

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob")

	mock.ExpectQuery("SELECT DISTINCT name FROM users").
		WillReturnRows(rows)

	names, err := service.GetUniqueNames()

	require.NoError(t, err)
	require.Equal(t, []string{"Alice", "Bob"}, names)
}

func TestGetUniqueNames_EmptyResult(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mock.ExpectationsWereMet()

	service := db.New(mockDB)

	rows := sqlmock.NewRows([]string{"name"})

	mock.ExpectQuery("SELECT DISTINCT name FROM users").
		WillReturnRows(rows)

	names, err := service.GetUniqueNames()

	require.NoError(t, err)
	require.Empty(t, names)
}

func TestGetUniqueNames_QueryError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mock.ExpectationsWereMet()

	service := db.New(mockDB)

	mock.ExpectQuery("SELECT DISTINCT name FROM users").
		WillReturnError(errors.New("db error"))

	names, err := service.GetUniqueNames()

	require.Error(t, err)
	require.Nil(t, names)
}

func TestGetUniqueNames_ScanError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mock.ExpectationsWereMet()

	service := db.New(mockDB)

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errors.New("scan error"))

	mock.ExpectQuery("SELECT DISTINCT name FROM users").
		WillReturnRows(rows)

	names, err := service.GetUniqueNames()

	require.Error(t, err)
	require.Nil(t, names)
}

func TestGetUniqueNames_RowsError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mock.ExpectationsWereMet()

	service := db.New(mockDB)

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		CloseError(errors.New("rows error"))

	mock.ExpectQuery("SELECT DISTINCT name FROM users").
		WillReturnRows(rows)

	names, err := service.GetUniqueNames()

	require.Error(t, err)
	require.Nil(t, names)
}
