package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ummmsh/task-6/internal/db"
)

var (
	errConnectionFailed = errors.New("connection failed")
	errCorruptedRow     = errors.New("corrupted row")
	errSyntaxError      = errors.New("syntax error")
	errIOFailure        = errors.New("io failure")
)

func TestNew(t *testing.T) {
	t.Parallel()

	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)

	defer mockDB.Close()

	service := db.New(mockDB)
	assert.Equal(t, mockDB, service.DB)
}

func TestDBServiceGetNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)

		defer mockDB.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob")

		mock.ExpectQuery("SELECT name FROM users").
			WillReturnRows(rows)

		service := db.New(mockDB)
		names, err := service.GetNames()

		require.NoError(t, err)
		assert.Equal(t, []string{"Alice", "Bob"}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)

		defer mockDB.Close()

		mock.ExpectQuery("SELECT name FROM users").
			WillReturnError(errConnectionFailed)

		service := db.New(mockDB)
		names, err := service.GetNames()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "db query")
		assert.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)

		defer mockDB.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("SELECT name FROM users").
			WillReturnRows(rows)

		service := db.New(mockDB)
		names, err := service.GetNames()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "rows scanning")
		assert.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows.Err() error", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)

		defer mockDB.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Test").
			RowError(0, errCorruptedRow)

		mock.ExpectQuery("SELECT name FROM users").
			WillReturnRows(rows)

		service := db.New(mockDB)
		names, err := service.GetNames()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "rows error")
		assert.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty result", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)

		defer mockDB.Close()

		rows := sqlmock.NewRows([]string{"name"})
		mock.ExpectQuery("SELECT name FROM users").
			WillReturnRows(rows)

		service := db.New(mockDB)
		names, err := service.GetNames()

		require.NoError(t, err)
		assert.Empty(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDBServiceGetUniqueNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)

		defer mockDB.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob")

		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnRows(rows)

		service := db.New(mockDB)
		names, err := service.GetUniqueNames()

		require.NoError(t, err)
		assert.Equal(t, []string{"Alice", "Bob"}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)

		defer mockDB.Close()

		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnError(errSyntaxError)

		service := db.New(mockDB)
		names, err := service.GetUniqueNames()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "db query")
		assert.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)

		defer mockDB.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnRows(rows)

		service := db.New(mockDB)
		names, err := service.GetUniqueNames()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "rows scanning")
		assert.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows.Err() error", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)

		defer mockDB.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Test").
			RowError(0, errIOFailure)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnRows(rows)

		service := db.New(mockDB)
		names, err := service.GetUniqueNames()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "rows error")
		assert.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty result", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)

		defer mockDB.Close()

		rows := sqlmock.NewRows([]string{"name"})
		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnRows(rows)

		service := db.New(mockDB)
		names, err := service.GetUniqueNames()

		require.NoError(t, err)
		assert.Empty(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
