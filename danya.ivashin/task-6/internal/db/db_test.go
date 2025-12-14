package db_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	database "github.com/Danya-byte/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

const (
	queryGetNames       = "SELECT name FROM users"
	queryGetUniqueNames = "SELECT DISTINCT name FROM users"
)

var ErrSome = errors.New("some error")

var casesGetNames = []struct { //nolint:gochecknoglobals
	names []string
}{
	{
		names: []string{"Ivan", "Gena228"},
	},
	{
		names: nil,
	},
}

func helperListMock(t *testing.T, values []string) *sqlmock.Rows {
	t.Helper()

	rows := sqlmock.NewRows([]string{"name"})
	for _, name := range values {
		rows = rows.AddRow(name)
	}

	return rows
}

func helperInitMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) { //nolint:ireturn
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	return db, mock
}

func TestGetNames_Success(t *testing.T) {
	t.Parallel()

	for _, test := range casesGetNames {
		db, mock := helperInitMock(t)
		defer db.Close()

		service := database.New(db)

		mock.ExpectQuery(queryGetNames).
			WillReturnRows(helperListMock(t, test.names))

		names, err := service.GetNames()

		require.Equal(t, test.names, names)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestGetNames_DBError(t *testing.T) {
	t.Parallel()

	db, mock := helperInitMock(t)
	defer db.Close()

	service := database.New(db)

	mock.ExpectQuery(queryGetNames).
		WillReturnRows(helperListMock(t, []string{"a", "b"})).
		WillReturnError(ErrSome)

	names, err := service.GetNames()

	require.Empty(t, names)
	require.ErrorIs(t, err, ErrSome)
	require.ErrorContains(t, err, "db query")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_RowsError(t *testing.T) {
	t.Parallel()

	db, mock := helperInitMock(t)
	defer db.Close()

	service := database.New(db)

	mock.ExpectQuery(queryGetNames).
		WillReturnRows(
			sqlmock.
				NewRows([]string{"name"}).
				AddRow("egor").
				RowError(0, ErrSome),
		)

	names, err := service.GetNames()

	require.ErrorIs(t, err, ErrSome)
	require.ErrorContains(t, err, "rows error")
	require.Empty(t, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_ScanError(t *testing.T) {
	t.Parallel()

	db, mock := helperInitMock(t)
	defer db.Close()

	service := database.New(db)

	mock.ExpectQuery(queryGetNames).
		WillReturnRows(
			sqlmock.
				NewRows([]string{"name"}).
				AddRow(nil),
		)

	names, err := service.GetNames()

	require.Error(t, err)
	require.ErrorContains(t, err, "rows scanning")
	require.Empty(t, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_Success(t *testing.T) {
	t.Parallel()

	for _, test := range casesGetNames {
		db, mock := helperInitMock(t)
		defer db.Close()

		service := database.New(db)

		mock.ExpectQuery(queryGetUniqueNames).
			WillReturnRows(helperListMock(t, test.names))

		names, err := service.GetUniqueNames()

		require.Equal(t, test.names, names)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestGetUniqueNames_DBError(t *testing.T) {
	t.Parallel()

	db, mock := helperInitMock(t)
	defer db.Close()

	service := database.New(db)

	mock.ExpectQuery(queryGetUniqueNames).
		WillReturnRows(helperListMock(t, []string{"a", "b"})).
		WillReturnError(ErrSome)

	names, err := service.GetUniqueNames()

	require.Empty(t, names)
	require.ErrorIs(t, err, ErrSome)
	require.ErrorContains(t, err, "db query")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_RowsError(t *testing.T) {
	t.Parallel()

	db, mock := helperInitMock(t)
	defer db.Close()

	service := database.New(db)

	mock.ExpectQuery(queryGetUniqueNames).
		WillReturnRows(
			sqlmock.
				NewRows([]string{"name"}).
				AddRow("egor").
				RowError(0, ErrSome),
		)

	names, err := service.GetUniqueNames()

	require.ErrorIs(t, err, ErrSome)
	require.ErrorContains(t, err, "rows error")
	require.Nil(t, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_ScanError(t *testing.T) {
	t.Parallel()

	db, mock := helperInitMock(t)
	defer db.Close()

	service := database.New(db)

	mock.ExpectQuery(queryGetUniqueNames).
		WillReturnRows(
			sqlmock.
				NewRows([]string{"name"}).
				AddRow(nil),
		)

	names, err := service.GetUniqueNames()

	require.Error(t, err)
	require.ErrorContains(t, err, "rows scanning")
	require.Empty(t, names)
	require.NoError(t, mock.ExpectationsWereMet())
}
