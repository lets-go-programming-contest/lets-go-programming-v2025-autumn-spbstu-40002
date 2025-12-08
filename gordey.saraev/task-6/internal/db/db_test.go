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
		defer DB.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob")

		mock.ExpectQuery(`^SELECT name FROM users$`).
			WillReturnRows(rows)

		names, err := db.New(DB).GetNames()
		require.NoError(t, err)
		require.Equal(t, []string{"Alice", "Bob"}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		DB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer DB.Close()

		mock.ExpectQuery(`^SELECT name FROM users$`).
			WillReturnError(errors.New("db connection failed"))

		_, err = db.New(DB).GetNames()
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		DB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer DB.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow(999) // неверный тип

		mock.ExpectQuery(`^SELECT name FROM users$`).
			WillReturnRows(rows)

		_, err = db.New(DB).GetNames()
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("row iteration error", func(t *testing.T) {
		DB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer DB.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			RowError(0, errors.New("network timeout"))

		mock.ExpectQuery(`^SELECT name FROM users$`).
			WillReturnRows(rows)

		_, err = db.New(DB).GetNames()
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows.Close() returns error", func(t *testing.T) {
		DB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer DB.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("John")
		rows.CloseError(errors.New("connection lost during close"))

		mock.ExpectQuery(`^SELECT name FROM users$`).
			WillReturnRows(rows)

		_, err = db.New(DB).GetNames()
		require.Error(t, err)
		require.Contains(t, err.Error(), "closing rows")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty result", func(t *testing.T) {
		DB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer DB.Close()

		rows := sqlmock.NewRows([]string{"name"})

		mock.ExpectQuery(`^SELECT name FROM users$`).
			WillReturnRows(rows)

		names, err := db.New(DB).GetNames()
		require.NoError(t, err)
		require.Empty(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Run("successful", func(t *testing.T) {
		DB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer DB.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob")

		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).
			WillReturnRows(rows)

			names, err := db.New(DB).GetUniqueNames()
		require.NoError(t, err)
		require.Equal(t, []string{"Alice", "Bob"}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows.Close() returns error", func(t *testing.T) {
		DB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer DB.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Unique")
		rows.CloseError(errors.New("close error")

		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).
			WillReturnRows(rows)

		_, err = db.New(DB).GetUniqueNames()
		require.Error(t, err)
		require.Contains(t, err.Error(), "closing rows")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		DB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer DB.Close()

		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).
			WillReturnError(errors.New("db error"))

		_, err = db.New(DB).GetUniqueNames()
		require.Error(t, err)
	})

	t.Run("scan error", func(t *testing.T) {
		DB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer DB.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)

		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).
			WillReturnRows(rows)

		_, err = db.New(DB).GetUniqueNames()
		require.Error(t, err)
	})

	t.Run("empty result", func(t *testing.T) {
		DB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer DB.Close()

		rows := sqlmock.NewRows([]string{"name"})

		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).
			WillReturnRows(rows)

		names, err := db.New(DB).GetUniqueNames()
		require.NoError(t, err)
		require.Empty(t, names)
	})
}