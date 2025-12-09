// internal/db/db_test.go
package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestDBService_GetNames(t *testing.T) {
	t.Run("successful", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").AddRow("Bob")
		mock.ExpectQuery(`^SELECT name FROM users$`).WillReturnRows(rows)

		names, err := db.New(db).GetNames()
		require.NoError(t, err)
		require.Equal(t, []string{"Alice", "Bob"}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(`^SELECT name FROM users$`).WillReturnError(errors.New("db down"))

		_, err = db.New(db).GetNames()
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow(123)
		mock.ExpectQuery(`^SELECT name FROM users$`).WillReturnRows(rows)

		_, err = db.New(db).GetNames()
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("row error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow("A").RowError(0, errors.New("boom"))
		mock.ExpectQuery(`^SELECT name FROM users$`).WillReturnRows(rows)

		_, err = db.New(db).GetNames()
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty result", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(`^SELECT name FROM users$`).WillReturnRows(sqlmock.NewRows([]string{"name"}))

		names, err := db.New(db).GetNames()
		require.NoError(t, err)
		require.Empty(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Run("successful", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").AddRow("Alice").AddRow("Bob")
		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).WillReturnRows(rows)

		names, err := db.New(db).GetUniqueNames()
		require.NoError(t, err)
		require.Equal(t, []string{"Alice", "Alice", "Bob"}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).WillReturnError(errors.New("db down"))

		_, err = db.New(db).GetUniqueNames()
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).WillReturnRows(rows)

		_, err = db.New(db).GetUniqueNames()
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty result", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).WillReturnRows(sqlmock.NewRows([]string{"name"}))

		names, err := db.New(db).GetUniqueNames()
		require.NoError(t, err)
		require.Empty(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
