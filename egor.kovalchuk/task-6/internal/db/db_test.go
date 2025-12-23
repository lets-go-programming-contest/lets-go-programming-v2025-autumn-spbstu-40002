package db_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/eg0sha-0/task-6/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	selectNamesQuery       = "SELECT name FROM users"
	selectUniqueNamesQuery = "SELECT DISTINCT name FROM users"
)

func TestDBService_GetNames_Success(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("John").
		AddRow("Gerald").
		AddRow("Michael")

	mock.ExpectQuery(selectNamesQuery).WillReturnRows(rows)

	service := db.New(mockDB)
	result, err := service.GetNames()

	require.NoError(t, err)
	assert.Equal(t, []string{"John", "Gerald", "Michael"}, result)
}

func TestDBService_GetNames_DBError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectQuery(selectNamesQuery).WillReturnError(sql.ErrConnDone)

	service := db.New(mockDB)
	result, err := service.GetNames()

	require.ErrorContains(t, err, "db query:")
	assert.Nil(t, result)
}

func TestDBService_GetNames_ScanFailure(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)

	mock.ExpectQuery(selectNamesQuery).WillReturnRows(rows)

	service := db.New(mockDB)
	result, err := service.GetNames()

	require.ErrorContains(t, err, "rows scanning:")
	assert.Nil(t, result)
}

func TestDBService_GetNames_RowIterationError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"name"})
	rows.AddRow("John")
	rows.RowError(0, sql.ErrTxDone)

	mock.ExpectQuery(selectNamesQuery).WillReturnRows(rows)

	service := db.New(mockDB)
	result, err := service.GetNames()

	require.ErrorContains(t, err, "rows error:")
	assert.Nil(t, result)
}

func TestDBService_GetNames_NoData(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"name"})
	mock.ExpectQuery(selectNamesQuery).WillReturnRows(rows)

	service := db.New(mockDB)
	result, err := service.GetNames()

	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestDBService_GetUniqueNames_Success(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("John").
		AddRow("Gerald").
		AddRow("John").
		AddRow("Michael").
		AddRow("Gerald")

	mock.ExpectQuery(selectUniqueNamesQuery).WillReturnRows(rows)

	service := db.New(mockDB)
	result, err := service.GetUniqueNames()

	require.NoError(t, err)
	assert.Equal(t, []string{"John", "Gerald", "John", "Michael", "Gerald"}, result)
}

func TestDBService_GetUniqueNames_DBError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectQuery(selectUniqueNamesQuery).WillReturnError(sql.ErrConnDone)

	service := db.New(mockDB)
	result, err := service.GetUniqueNames()

	require.ErrorContains(t, err, "db query:")
	assert.Nil(t, result)
}

func TestDBService_GetUniqueNames_ScanFailure(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)

	mock.ExpectQuery(selectUniqueNamesQuery).WillReturnRows(rows)

	service := db.New(mockDB)
	result, err := service.GetUniqueNames()

	require.ErrorContains(t, err, "rows scanning:")
	assert.Nil(t, result)
}

func TestDBService_GetUniqueNames_RowIterationError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"name"})
	rows.AddRow("John")
	rows.RowError(0, sql.ErrTxDone)

	mock.ExpectQuery(selectUniqueNamesQuery).WillReturnRows(rows)

	service := db.New(mockDB)
	result, err := service.GetUniqueNames()

	require.ErrorContains(t, err, "rows error:")
	assert.Nil(t, result)
}

func TestDBService_GetUniqueNames_NoData(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"name"})
	mock.ExpectQuery(selectUniqueNamesQuery).WillReturnRows(rows)

	service := db.New(mockDB)
	result, err := service.GetUniqueNames()

	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestDBService_Initialization(t *testing.T) {
	t.Parallel()

	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)

	assert.NotNil(t, service)
	assert.Equal(t, mockDB, service.DB)
}
