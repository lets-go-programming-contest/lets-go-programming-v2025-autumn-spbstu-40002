package db_test

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mydb "github.com/stepanov.alexander/task-6/internal/db"
)

func TestDBService_GetNames(t *testing.T) {
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

		service := mydb.New(dbMock)

		names, err := service.GetNames()
		require.NoError(t, err)
		assert.NotEmpty(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success empty", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		rows := sqlmock.NewRows([]string{"name"})

		mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).WillReturnRows(rows)

		service := mydb.New(dbMock)

		names, err := service.GetNames()
		require.NoError(t, err)
		assert.Empty(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).WillReturnError(assert.AnError)

		service := mydb.New(dbMock)

		names, err := service.GetNames()
		require.Error(t, err)
		require.ErrorContains(t, err, "db query")
		assert.Nil(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).WillReturnRows(rows)

		service := mydb.New(dbMock)

		names, err := service.GetNames()
		require.Error(t, err)
		require.ErrorContains(t, err, "rows scanning")
		assert.Nil(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			RowError(0, assert.AnError)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).WillReturnRows(rows)

		service := mydb.New(dbMock)

		names, err := service.GetNames()
		require.Error(t, err)
		require.ErrorContains(t, err, "rows error")
		assert.Nil(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
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

		service := mydb.New(dbMock)

		names, err := service.GetUniqueNames()
		require.NoError(t, err)
		assert.NotEmpty(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success empty", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		rows := sqlmock.NewRows([]string{"name"})

		mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).WillReturnRows(rows)

		service := mydb.New(dbMock)

		names, err := service.GetUniqueNames()
		require.NoError(t, err)
		assert.Empty(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).WillReturnError(assert.AnError)

		service := mydb.New(dbMock)

		names, err := service.GetUniqueNames()
		require.Error(t, err)
		require.ErrorContains(t, err, "db query")
		assert.Nil(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).WillReturnRows(rows)

		service := mydb.New(dbMock)

		names, err := service.GetUniqueNames()
		require.Error(t, err)
		require.ErrorContains(t, err, "rows scanning")
		assert.Nil(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			RowError(0, assert.AnError)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).WillReturnRows(rows)

		service := mydb.New(dbMock)

		names, err := service.GetUniqueNames()
		require.Error(t, err)
		require.ErrorContains(t, err, "rows error")
		assert.Nil(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDBService_New(t *testing.T) {
	t.Parallel()

	dbMock, _, err := sqlmock.New()
	require.NoError(t, err)
	defer dbMock.Close()

	service := mydb.New(dbMock)
	assert.NotNil(t, service)
}
