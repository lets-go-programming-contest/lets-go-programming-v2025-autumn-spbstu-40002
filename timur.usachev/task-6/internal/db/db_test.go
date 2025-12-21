package db_test

import (
	"errors"
	"regexp"
	"testing"

	db "github.com/t1wt/task-6/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

var (
	errQuery = errors.New("boom")
	errRows  = errors.New("rows error")
)

func closeDB(t *testing.T, c interface{ Close() error }) {
	t.Helper()

	if err := c.Close(); err != nil {
		t.Logf("close db: %v", err)
	}
}

func TestGetNames_SuccessWithData(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer closeDB(t, sqlDB)

	service := db.New(sqlDB)
	rows := sqlmock.NewRows([]string{"name"}).AddRow("alice").AddRow("bob")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).WillReturnRows(rows)

	names, err := service.GetNames()
	require.NoError(t, err)
	require.Equal(t, []string{"alice", "bob"}, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_SuccessNoData(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer closeDB(t, sqlDB)

	service := db.New(sqlDB)
	rows := sqlmock.NewRows([]string{"name"})
	mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).WillReturnRows(rows)

	names, err := service.GetNames()
	require.NoError(t, err)
	require.Empty(t, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_QueryError(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer closeDB(t, sqlDB)

	service := db.New(sqlDB)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).WillReturnError(errQuery)

	names, err := service.GetNames()
	require.Nil(t, names)
	require.ErrorIs(t, err, errQuery)
	require.Regexp(t, "^db query: .*boom", err.Error())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_ScanError(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer closeDB(t, sqlDB)

	service := db.New(sqlDB)
	rows := sqlmock.NewRows([]string{"name", "age"}).AddRow("alice", 30)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).WillReturnRows(rows)

	names, err := service.GetNames()
	require.Nil(t, names)
	require.Regexp(t, "rows scanning: .*", err.Error())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_RowsError(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer closeDB(t, sqlDB)

	service := db.New(sqlDB)
	rows := sqlmock.NewRows([]string{"name"}).AddRow("alice").RowError(0, errRows)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).WillReturnRows(rows)

	names, err := service.GetNames()
	require.Nil(t, names)
	require.Regexp(t, "rows error: .*rows error", err.Error())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_SuccessWithData(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer closeDB(t, sqlDB)

	service := db.New(sqlDB)
	rows := sqlmock.NewRows([]string{"name"}).AddRow("alice").AddRow("alice")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).WillReturnRows(rows)

	names, err := service.GetUniqueNames()
	require.NoError(t, err)
	require.Equal(t, []string{"alice", "alice"}, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_SuccessNoData(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer closeDB(t, sqlDB)

	service := db.New(sqlDB)
	rows := sqlmock.NewRows([]string{"name"})
	mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).WillReturnRows(rows)

	names, err := service.GetUniqueNames()
	require.NoError(t, err)
	require.Empty(t, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_QueryError(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer closeDB(t, sqlDB)

	service := db.New(sqlDB)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).WillReturnError(errQuery)

	names, err := service.GetUniqueNames()
	require.Nil(t, names)
	require.ErrorIs(t, err, errQuery)
	require.Regexp(t, "^db query: .*boom", err.Error())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_ScanError(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer closeDB(t, sqlDB)

	service := db.New(sqlDB)
	rows := sqlmock.NewRows([]string{"name", "age"}).AddRow("alice", 30)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).WillReturnRows(rows)

	names, err := service.GetUniqueNames()
	require.Nil(t, names)
	require.Regexp(t, "rows scanning: .*", err.Error())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_RowsError(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer closeDB(t, sqlDB)

	service := db.New(sqlDB)
	rows := sqlmock.NewRows([]string{"name"}).AddRow("alice").RowError(0, errRows)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).WillReturnRows(rows)

	names, err := service.GetUniqueNames()
	require.Nil(t, names)
	require.Regexp(t, "rows error: .*rows error", err.Error())
	require.NoError(t, mock.ExpectationsWereMet())
}
