package db_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/F0LY/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

type errorDatabase struct {
	err error
}

func (e *errorDatabase) Query(query string, args ...any) (*sql.Rows, error) {
	return nil, e.err
}

func TestNewDBService(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	service := db.New(sqlDB)
	require.NotNil(t, service)
	require.NotNil(t, service.DB)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetNames(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	service := db.New(sqlDB)

	t.Run("successful", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").AddRow("Bob")
		mock.ExpectQuery(`^SELECT name FROM users$`).WillReturnRows(rows)

		names, err := service.GetNames()
		require.NoError(t, err)
		require.Equal(t, []string{"Alice", "Bob"}, names)
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT name FROM users$`).WillReturnError(errors.New("db error"))
		_, err := service.GetNames()
		require.Error(t, err)
	})

	t.Run("scan error type mismatch", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow(123)
		mock.ExpectQuery(`^SELECT name FROM users$`).WillReturnRows(rows)
		_, err := service.GetNames()
		require.Error(t, err)
	})

	t.Run("scan null value error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery(`^SELECT name FROM users$`).WillReturnRows(rows)
		names, err := service.GetNames()
		require.Error(t, err)
		require.Contains(t, err.Error(), "rows scanning:")
		require.Nil(t, names)
	})

	t.Run("row error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow("A").RowError(0, errors.New("boom"))
		mock.ExpectQuery(`^SELECT name FROM users$`).WillReturnRows(rows)
		_, err := service.GetNames()
		require.Error(t, err)
	})

	t.Run("rows.Err() error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow("test")
		rows.RowError(1, errors.New("delayed error"))
		mock.ExpectQuery(`^SELECT name FROM users$`).WillReturnRows(rows)
		_, err := service.GetNames()
		require.Error(t, err)
		require.Contains(t, err.Error(), "rows error")
	})

	t.Run("rows close error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").CloseError(errors.New("close error"))
		mock.ExpectQuery(`^SELECT name FROM users$`).WillReturnRows(rows)
		names, err := service.GetNames()
		require.Error(t, err)
		require.Contains(t, err.Error(), "rows error:")
		require.Nil(t, names)
	})

	t.Run("empty result", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT name FROM users$`).WillReturnRows(sqlmock.NewRows([]string{"name"}))
		names, err := service.GetNames()
		require.NoError(t, err)
		require.Empty(t, names)
	})

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetUniqueNames(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	service := db.New(sqlDB)

	t.Run("successful", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").AddRow("Bob")
		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).WillReturnRows(rows)
		names, err := service.GetUniqueNames()
		require.NoError(t, err)
		require.Equal(t, []string{"Alice", "Bob"}, names)
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).WillReturnError(errors.New("db error"))
		_, err := service.GetUniqueNames()
		require.Error(t, err)
	})

	t.Run("scan error type mismatch", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow(123)
		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).WillReturnRows(rows)
		_, err := service.GetUniqueNames()
		require.Error(t, err)
	})

	t.Run("scan null value error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).WillReturnRows(rows)
		names, err := service.GetUniqueNames()
		require.Error(t, err)
		require.Contains(t, err.Error(), "rows scanning:")
		require.Nil(t, names)
	})

	t.Run("rows.Err() error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow("test")
		rows.RowError(1, errors.New("delayed error"))
		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).WillReturnRows(rows)
		_, err := service.GetUniqueNames()
		require.Error(t, err)
		require.Contains(t, err.Error(), "rows error")
	})

	t.Run("rows close error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").CloseError(errors.New("close error"))
		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).WillReturnRows(rows)
		names, err := service.GetUniqueNames()
		require.Error(t, err)
		require.Contains(t, err.Error(), "rows error:")
		require.Nil(t, names)
	})

	t.Run("empty result", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).WillReturnRows(sqlmock.NewRows([]string{"name"}))
		names, err := service.GetUniqueNames()
		require.NoError(t, err)
		require.Empty(t, names)
	})

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_WithErrorDatabase(t *testing.T) {
	t.Run("GetNames with custom error database", func(t *testing.T) {
		mockErr := errors.New("custom query error")
		dbService := db.New(&errorDatabase{err: mockErr})

		names, err := dbService.GetNames()
		require.Error(t, err)
		require.Contains(t, err.Error(), "db query:")
		require.Nil(t, names)
	})

	t.Run("GetUniqueNames with custom error database", func(t *testing.T) {
		mockErr := errors.New("custom query error")
		dbService := db.New(&errorDatabase{err: mockErr})

		names, err := dbService.GetUniqueNames()
		require.Error(t, err)
		require.Contains(t, err.Error(), "db query:")
		require.Nil(t, names)
	})
}
