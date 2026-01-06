package db

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	return db, mock
}

func TestDBService_GetNames(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alexander").
		AddRow("Petr").
		AddRow("Feodosiy")

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := DBService{DB: db}

	names, err := service.GetNames()

	require.NoError(t, err)
	require.NotNil(t, names)
	assert.Equal(t, []string{"Alexander", "Petr", "Feodosiy"}, names)

	require.NoError(t, mock.ExpectationsWereMet())

}
