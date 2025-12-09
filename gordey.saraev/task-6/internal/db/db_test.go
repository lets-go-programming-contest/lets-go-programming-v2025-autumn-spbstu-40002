// internal/db/db_test.go
package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestDBService_GetNames(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	t.Run("ok", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow("A").AddRow("B")
		mock.ExpectQuery("^SELECT name FROM users$").WillReturnRows(rows)
		names, err := db.New(db).GetNames()
		require.NoError(t, err)
		require.Equal(t, []string{"A", "B"}, names)
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery("^SELECT name FROM users$").WillReturnError(errors.New("boom"))
		_, err := db.New(db).GetNames()
		require.Error(t, err)
	})

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetUniqueNames(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	t.Run("ok", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow("A").AddRow("B")
		mock.ExpectQuery("^SELECT DISTINCT name FROM users$").WillReturnRows(rows)
		names, err := db.New(db).GetUniqueNames()
		require.NoError(t, err)
		require.Equal(t, []string{"A", "B"}, names)
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery("^SELECT DISTINCT name FROM users$").WillReturnError(errors.New("boom"))
		_, err := db.New(db).GetUniqueNames()
		require.Error(t, err)
	})

	require.NoError(t, mock.ExpectationsWereMet())
}
