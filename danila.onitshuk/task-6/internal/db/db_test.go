package db_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/danila-clown/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

var (
	errQueryFailed      = errors.New("query failed")
	errRowProcessing    = errors.New("row processing failed")
	errDatabaseAccess   = errors.New("database access failed")
	errResultProcessing = errors.New("result processing failed")
	errOperationFailed  = errors.New("operation failed")
)

type errorDatabase struct {
	err error
}

func (e *errorDatabase) Query(query string, args ...any) (*sql.Rows, error) {
	return nil, e.err
}

func TestNewDBService(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	service := db.New(sqlDB)
	require.NotNil(t, service)
	require.NotNil(t, service.DB)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	t.Run("successful", func(t *testing.T) {
		t.Parallel()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").AddRow("Bob")
		mock.ExpectQuery(`^SELECT name FROM users$`).WillReturnRows(rows)

		names, err := service.GetNames()
		require.NoError(t, err)
		require.Equal(t, []string{"Alice", "Bob"}, names)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)

		mock.ExpectQuery(`^SELECT name FROM users$`).WillReturnError(errQueryFailed)

		_, err = service.GetNames()
		require.Error(t, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan null value error", func(t *testing.T) {
		t.Parallel()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery(`^SELECT name FROM users$`).WillReturnRows(rows)

		names, err := service.GetNames()
		require.Error(t, err)
		require.Contains(t, err.Error(), "rows scanning:")
		require.Nil(t, names)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("row error", func(t *testing.T) {
		t.Parallel()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow("A").RowError(0, errRowProcessing)
		mock.ExpectQuery(`^SELECT name FROM users$`).WillReturnRows(rows)

		_, err = service.GetNames()
		require.Error(t, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows close error", func(t *testing.T) {
		t.Parallel()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").CloseError(errResultProcessing)
		mock.ExpectQuery(`^SELECT name FROM users$`).WillReturnRows(rows)

		names, err := service.GetNames()
		require.Error(t, err)
		require.Contains(t, err.Error(), "rows error:")
		require.Nil(t, names)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty result", func(t *testing.T) {
		t.Parallel()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)

		mock.ExpectQuery(`^SELECT name FROM users$`).WillReturnRows(sqlmock.NewRows([]string{"name"}))

		names, err := service.GetNames()
		require.NoError(t, err)
		require.Empty(t, names)

		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	t.Run("successful", func(t *testing.T) {
		t.Parallel()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").AddRow("Bob")
		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).WillReturnRows(rows)

		names, err := service.GetUniqueNames()
		require.NoError(t, err)
		require.Equal(t, []string{"Alice", "Bob"}, names)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)

		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).WillReturnError(errDatabaseAccess)

		_, err = service.GetUniqueNames()
		require.Error(t, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan null value error", func(t *testing.T) {
		t.Parallel()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).WillReturnRows(rows)

		names, err := service.GetUniqueNames()
		require.Error(t, err)
		require.Contains(t, err.Error(), "rows scanning:")
		require.Nil(t, names)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows close error", func(t *testing.T) {
		t.Parallel()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").CloseError(errResultProcessing)
		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).WillReturnRows(rows)

		names, err := service.GetUniqueNames()
		require.Error(t, err)
		require.Contains(t, err.Error(), "rows error:")
		require.Nil(t, names)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty result", func(t *testing.T) {
		t.Parallel()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)

		mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).WillReturnRows(sqlmock.NewRows([]string{"name"}))

		names, err := service.GetUniqueNames()
		require.NoError(t, err)
		require.Empty(t, names)

		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDBService_WithErrorDatabase(t *testing.T) {
	t.Parallel()

	t.Run("GetNames with custom error database", func(t *testing.T) {
		t.Parallel()

		dbService := db.New(&errorDatabase{err: errOperationFailed})

		names, err := dbService.GetNames()
		require.Error(t, err)
		require.Contains(t, err.Error(), "db query:")
		require.Nil(t, names)
	})

	t.Run("GetUniqueNames with custom error database", func(t *testing.T) {
		t.Parallel()

		dbService := db.New(&errorDatabase{err: errOperationFailed})

		names, err := dbService.GetUniqueNames()
		require.Error(t, err)
		require.Contains(t, err.Error(), "db query:")
		require.Nil(t, names)
	})
}
