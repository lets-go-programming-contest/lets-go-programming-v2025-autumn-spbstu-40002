package db_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/leonid.maryankov/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

var (
	errDB           = errors.New("db error")
	errRowIteration = errors.New("row iteration error")
)

func TestNew(t *testing.T) {
	t.Parallel()

	dbMock, _, err := sqlmock.New()
	require.NoError(t, err)
	defer dbMock.Close()

	service := db.New(dbMock)
	require.NotNil(t, service)
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	t.Run("success with data", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).WillReturnRows(rows)

		service := db.New(dbMock)
		names, err := service.GetNames()

		require.NoError(t, err)
		require.Equal(t, []string{"Alice", "Bob"}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success empty", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		rows := sqlmock.NewRows([]string{"name"})
		mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).WillReturnRows(rows)

		service := db.New(dbMock)
		names, err := service.GetNames()

		require.NoError(t, err)
		require.Empty(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).WillReturnError(errDB)

		service := db.New(dbMock)
		names, err := service.GetNames()

		require.Error(t, err)
		require.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).WillReturnRows(rows)

		service := db.New(dbMock)
		names, err := service.GetNames()

		require.Error(t, err)
		require.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			RowError(0, errRowIteration)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).WillReturnRows(rows)

		service := db.New(dbMock)
		names, err := service.GetNames()

		require.Error(t, err)
		require.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetUniqueNames(t *testing.T) {
	t.Parallel()

	t.Run("success with data", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).WillReturnRows(rows)

		service := db.New(dbMock)
		names, err := service.GetUniqueNames()

		require.NoError(t, err)
		require.Equal(t, []string{"Alice", "Bob"}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success empty", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		rows := sqlmock.NewRows([]string{"name"})
		mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).WillReturnRows(rows)

		service := db.New(dbMock)
		names, err := service.GetUniqueNames()

		require.NoError(t, err)
		require.Empty(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).WillReturnError(errDB)

		service := db.New(dbMock)
		names, err := service.GetUniqueNames()

		require.Error(t, err)
		require.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).WillReturnRows(rows)

		service := db.New(dbMock)
		names, err := service.GetUniqueNames()

		require.Error(t, err)
		require.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			RowError(0, errRowIteration)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).WillReturnRows(rows)

		service := db.New(dbMock)
		names, err := service.GetUniqueNames()

		require.Error(t, err)
		require.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
