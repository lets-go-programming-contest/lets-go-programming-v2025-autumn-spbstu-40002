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
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob")

		mock.ExpectQuery(`^SELECT name FROM users$`).
			WillReturnRows(rows)

		names, err := service.GetNames()
		require.NoError(t, err)
		require.Equal(t, []string{"Alice", "Bob"}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT name FROM users$`).
			WillReturnError(errors.New("db error"))

		_, err := service.GetNames()
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow(123)

		mock.ExpectQuery(`^SELECT name FROM users$`).
			WillReturnRows(rows)

		_, err := service.GetNames()
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("row error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			RowError(0, errors.New("network error"))

		mock.ExpectQuery(`^SELECT name FROM users$`).
			WillReturnRows(rows)

		_, err := service.GetNames()
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows.Err() error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("test").
			RowError(1, errors.New("delayed error"))

		mock.ExpectQuery(`^SELECT name FROM users$`).
			WillReturnRows(rows)

		_, err := service.GetNames()
		require.Error(t, err)
		require.Contains(t, err.Error(), "rows error")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty result", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT name FROM users$`).
			WillReturnRows(sqlmock.NewRows([]string{"name"}))

		names, err := service.GetNames()
		require.NoError(t, err)
		require.Empty(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	service := db.New(sqlDB)

	t.Run("successful", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob")

		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).
			WillReturnRows(rows)

		names, err := service.GetUniqueNames()
		require.NoError(t, err)
		require.Equal(t, []string{"Alice", "Bob"}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).
			WillReturnError(errors.New("db error"))

		_, err := service.GetUniqueNames()
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)

		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).
			WillReturnRows(rows)

		_, err := service.GetUniqueNames()
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows.Err() error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("test").
			RowError(1, errors.New("delayed error"))

		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).
			WillReturnRows(rows)

		_, err := service.GetUniqueNames()
		require.Error(t, err)
		require.Contains(t, err.Error(), "rows error")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty result", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).
			WillReturnRows(sqlmock.NewRows([]string{"name"}))

		names, err := service.GetUniqueNames()
		require.NoError(t, err)
		require.Empty(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
