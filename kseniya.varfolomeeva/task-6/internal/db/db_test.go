package db_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	Mdb "kseniya.varfolomeeva/task-6/internal/db"
)

func TestNew(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectationsWereMet()

	service := Mdb.New(db)
	require.Equal(t, db, service.DB, "Expected DB to be set")
}

func TestDBService_GetNames_Success(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Ivan").
		AddRow("Maria").
		AddRow("Alexey").
		AddRow("Eugene").
		AddRow("Olga")
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := Mdb.New(db)
	names, err := service.GetNames()
	require.NoError(t, err)

	require.Len(t, names, 5, "Expected 5 names")
	require.Equal(t, "Ivan", names[0], "First name should be Ivan")
	require.Equal(t, "Maria", names[1], "Second name should be Maria")
	require.Equal(t, "Alexey", names[2], "Third name should be Alexey")
	require.Equal(t, "Eugene", names[3], "Fourth name should be Eugene")
	require.Equal(t, "Olga", names[4], "Fifth name should be Olga")

	require.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestDBService_GetNames_WithDuplicates(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Ivan").
		AddRow("Maria").
		AddRow("Ivan").
		AddRow("Eugene").
		AddRow("Maria")
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := Mdb.New(db)
	names, err := service.GetNames()
	require.NoError(t, err)

	require.Len(t, names, 5, "Expected 5 names with duplicates")
	require.Equal(t, []string{"Ivan", "Maria", "Ivan", "Eugene", "Maria"}, names)

	require.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestDBService_GetNames_Empty(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"})
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := Mdb.New(db)
	names, err := service.GetNames()
	require.NoError(t, err)
	require.Empty(t, names, "Expected empty slice")

	require.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestDBService_GetNames_QueryError(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnError(sql.ErrConnDone)

	service := Mdb.New(db)
	names, err := service.GetNames()
	require.Error(t, err, "Expected error")
	require.Nil(t, names, "Expected nil result on error")
	require.Contains(t, err.Error(), "db query", "Error should contain 'db query'")

	require.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestDBService_GetNames_ScanError(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Ivan").
		AddRow(nil)
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := Mdb.New(db)
	names, err := service.GetNames()
	require.Error(t, err, "Expected error")
	require.Nil(t, names, "Expected nil result on error")
	require.Contains(t, err.Error(), "rows scanning", "Error should contain 'rows scanning'")

	require.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestDBService_GetNames_RowsError(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"})
	rows.AddRow("Ivan")
	rows.RowError(0, sql.ErrTxDone)
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := Mdb.New(db)
	names, err := service.GetNames()
	require.Error(t, err, "Expected error")
	require.Nil(t, names, "Expected nil result on error")
	require.Contains(t, err.Error(), "rows error", "Error should contain 'rows error'")

	require.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestDBService_GetUniqueNames_Success(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Ivan").
		AddRow("Maria").
		AddRow("Alexey").
		AddRow("Eugene").
		AddRow("Olga").
		AddRow("Ivan").
		AddRow("Eugene")
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	service := Mdb.New(db)
	names, err := service.GetUniqueNames()
	require.NoError(t, err)
	require.Len(t, names, 7, "Expected 7 names (with duplicates from DISTINCT)")

	expected := []string{"Ivan", "Maria", "Alexey", "Eugene", "Olga", "Ivan", "Eugene"}
	for i, name := range names {
		require.Equal(t, expected[i], name, "Name mismatch at index %d", i)
	}

	require.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestDBService_GetUniqueNames_OnlyEugene(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Eugene").
		AddRow("Eugene").
		AddRow("Eugene")
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	service := Mdb.New(db)
	names, err := service.GetUniqueNames()
	require.NoError(t, err)
	require.Len(t, names, 3, "Expected 3 Eugene names")

	for i, name := range names {
		require.Equal(t, "Eugene", name, "All names should be Eugene at index %d", i)
	}

	require.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestDBService_GetUniqueNames_Empty(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"})
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	service := Mdb.New(db)
	names, err := service.GetUniqueNames()
	require.NoError(t, err)
	require.Empty(t, names, "Expected empty slice")

	require.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestDBService_GetUniqueNames_QueryError(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT DISTINCT name FROM users").
		WillReturnError(sql.ErrConnDone)

	service := Mdb.New(db)
	names, err := service.GetUniqueNames()
	require.Error(t, err, "Expected error")
	require.Nil(t, names, "Expected nil result on error")
	require.Contains(t, err.Error(), "db query", "Error should contain 'db query'")

	require.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestDBService_GetUniqueNames_ScanError(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Ivan").
		AddRow(nil)
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	service := Mdb.New(db)
	names, err := service.GetUniqueNames()
	require.Error(t, err, "Expected error")
	require.Nil(t, names, "Expected nil result on error")
	require.Contains(t, err.Error(), "rows scanning", "Error should contain 'rows scanning'")

	require.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestDBService_GetUniqueNames_RowsError(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"})
	rows.AddRow("Ivan")
	rows.RowError(0, sql.ErrTxDone)
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	service := Mdb.New(db)
	names, err := service.GetUniqueNames()
	require.Error(t, err, "Expected error")
	require.Nil(t, names, "Expected nil result on error")
	require.Contains(t, err.Error(), "rows error", "Error should contain 'rows error'")

	require.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestDBService_GetNames_QueryErrorClosesRows(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnError(sql.ErrConnDone)

	service := Mdb.New(db)
	names, err := service.GetNames()
	require.Error(t, err)
	require.Nil(t, names)
	require.Contains(t, err.Error(), "db query")

	require.NoError(t, mock.ExpectationsWereMet())
}
