package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/F0LY/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

func TestDBService_GetNames(t *testing.T) {
	t.Run("successful", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob")

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := service.GetNames()

		require.NoError(t, err)
		require.Equal(t, []string{"Alice", "Bob"}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)

		mock.ExpectQuery("SELECT name FROM users").WillReturnError(errors.New("query error"))

		_, err = service.GetNames()

		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow(123) // type mismatch

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		_, err = service.GetNames()

		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows err", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice")

		rows.RowError(0, errors.New("row error"))

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		_, err = service.GetNames()

		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"})

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := service.GetNames()

		require.NoError(t, err)
		require.Empty(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Run("successful", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob")

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := service.GetUniqueNames()

		require.NoError(t, err)
		require.Equal(t, []string{"Alice", "Bob"}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errors.New("query error"))

		_, err = service.GetUniqueNames()

		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow(123) // type mismatch

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		_, err = service.GetUniqueNames()

		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows err", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice")

		rows.RowError(0, errors.New("row error"))

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		_, err = service.GetUniqueNames()

		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"})

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := service.GetUniqueNames()

		require.NoError(t, err)
		require.Empty(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
