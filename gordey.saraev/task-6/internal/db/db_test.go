package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/F0LY/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

func TestDBService_GetNames(t *testing.T) {
	t.Run("successful — covers rows.Close() via normal execution", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)

		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob")

		mock.ExpectQuery(`^SELECT name FROM users$`).
			WillReturnRows(rows)

		names, err := service.GetNames()
		require.NoError(t, err)
		require.Equal(t, []string{"Alice", "Bob"}, names)
		require.NoError(t, mock.ExpectationsWereMet())

		sqlDB.Close()
	})

	t.Run("query error", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		mock.ExpectQuery(`^SELECT name FROM users$`).
			WillReturnError(errors.New("db down"))

		_, err = db.New(sqlDB).GetNames()
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow(123)

		mock.ExpectQuery(`^SELECT name FROM users$`).
			WillReturnRows(rows)

		_, err = db.New(sqlDB).GetNames()
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("row error during iteration", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			RowError(0, errors.New("timeout"))

		mock.ExpectQuery(`^SELECT name FROM users$`).
			WillReturnRows(rows)

		_, err = db.New(sqlDB).GetNames()
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty result", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		rows := sqlmock.NewRows([]string{"name"})

		mock.ExpectQuery(`^SELECT name FROM users$`).
			WillReturnRows(rows)

		names, err := db.New(sqlDB).GetNames()
		require.NoError(t, err)
		require.Empty(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Run("successful — also covers rows.Close() via normal execution", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)

		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob")

		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).
			WillReturnRows(rows)

		names, err := service.GetUniqueNames()
		require.NoError(t, err)
		require.Equal(t, []string{"Alice", "Bob"}, names)
		require.NoError(t, mock.ExpectationsWereMet())

		sqlDB.Close()
	})

}
