package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/F0LY/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

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
	)
		require.Equal(t, []string{"Alice", "Bob"}, names)
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT name FROM users$`).WillReturnError(errors.New("db error"))
		_, err := service.GetNames()
		require.Error(t, err)
	})

	t.Run("scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow(123)
		mock.ExpectQuery(`^SELECT name FROM users$`).WillReturnRows(rows)

		_, err := service.GetNames()
		require.Error(t, err)
	})

	t.Run("row error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow("A").RowError(0, errors.New("boom"))
		mock.ExpectQuery(`^SELECT name FROM users$`).WillReturnRows(rows)

		_, err := service.GetNames()
		require.Error(t, err)
	})

	t.Run("rows.Close() error", func(t *testing.T) {
		sqlDB2, mock2, err := sqlmock.New()
		require.NoError(t, err)

		rows := sqlmock.NewRows([]string{"name"}).AddRow("test")
		rows.CloseError(errors.New("connection lost"))

		mock2.ExpectQuery(`^SELECT name FROM users$`).WillReturnRows(rows)

		service2 := db.New(sqlDB2)
		_, err = service2.GetNames()
		require.Error(t, err)
		require.Contains(t, err.Error(), "closing rows")
		require.NoError(t, mock2.ExpectationsWereMet())

		sqlDB2.Close()
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

	t.Run("scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).WillReturnRows(rows)

		_, err := service.GetUniqueNames()
		require.Error(t, err)
	})

	t.Run("rows.Close() error", func(t *testing.T) {
		sqlDB2, mock2, err := sqlmock.New()
		require.NoError(t, err)

		rows := sqlmock.NewRows([]string{"name"}).AddRow("test")
		rows.CloseError(errors.New("connection lost"))

		mock2.ExpectQuery(`^SELECT DISTINCT name FROM users$`).WillReturnRows(rows)

		service2 := db.New(sqlDB2)
		_, err = service2.GetUniqueNames()
		require.Error(t, err)
		require.Contains(t, err.Error(), "closing rows")
		require.NoError(t, mock2.ExpectationsWereMet())

		sqlDB2.Close()
	})

	require.NoError(t, mock.ExpectationsWereMet())
}