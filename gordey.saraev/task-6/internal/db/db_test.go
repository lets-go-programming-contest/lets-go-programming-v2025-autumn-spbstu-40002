// internal/db/db_test.go
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

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob")

		mock.ExpectQuery(`^SELECT name FROM users$`).
			WillReturnRows(rows)

		names, err := db.New(sqlDB).GetNames()
		require.NoError(t, err)
		require.Equal(t, []string{"Alice", "Bob"}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		mock.ExpectQuery(`^SELECT name FROM users$`).
			WillReturnError(errors.New("db error"))

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

	t.Run("row error", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			RowError(0, errors.New("boom"))

		mock.ExpectQuery(`^SELECT name FROM users$`).
			WillReturnRows(rows)

		_, err = db.New(sqlDB).GetNames()
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows.Err() error", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("test")
		rows.RowError(1, errors.New("delayed"))

		mock.ExpectQuery(`^SELECT name FROM users$`).
			WillReturnRows(rows)

		_, err = db.New(sqlDB).GetNames()
		require.Error(t, err)
		require.Contains(t, err.Error(), "rows error")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		mock.ExpectQuery(`^SELECT name FROM users$`).
			WillReturnRows(sqlmock.NewRows([]string{"name"}))

		names, err := db.New(sqlDB).GetNames()
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

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob")

		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).
			WillReturnRows(rows)

		names, err := db.New(sqlDB).GetUniqueNames()
		require.NoError(t, err)
		require.Equal(t, []string{"Alice", "Bob"}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).
			WillReturnError(errors.New("db error"))

		_, err = db.New(sqlDB).GetUniqueNames()
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)

		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).
			WillReturnRows(rows)

		_, err = db.New(sqlDB(sqlDB).GetUniqueNames()
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows.Err() error", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("test")
		rows.RowError(1, errors.New("delayed"))

		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).
			WillReturnRows(rows)

		_, err = db.New(sqlDB).GetUniqueNames()
		require.Error(t, err)
		require.Contains(t, err.Error(), "rows error")
		require.NoError(t, mock.ExpectationsWereMet())
	})
}