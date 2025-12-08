package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDBService_GetNames(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	service := New(db)

	tests := []struct {
		name         string
		mockNames    []string
		mockErr      error
		scanErr      bool
		rowsErr      bool
		expectErr    bool
		expectResult []string
	}{
		{
			name:         "success with multiple names",
			mockNames:    []string{"Alice", "Bob", "Charlie"},
			expectResult: []string{"Alice", "Bob", "Charlie"},
		},
		{
			name:         "success with empty result",
			mockNames:    []string{},
			expectResult: []string{},
		},
		{
			name:      "query error",
			mockErr:   errors.New("query failed"),
			expectErr: true,
		},
		{
			name:         "single name",
			mockNames:    []string{"John"},
			expectResult: []string{"John"},
		},
		{
			name:      "scan error",
			scanErr:   true,
			expectErr: true,
		},
		{
			name:      "rows error",
			rowsErr:   true,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := "SELECT name FROM users"

			if tt.mockErr != nil {
				mock.ExpectQuery(query).WillReturnError(tt.mockErr)
			} else if tt.scanErr {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice")
				rows.RowError(0, errors.New("scan error"))
				mock.ExpectQuery(query).WillReturnRows(rows)
			} else if tt.rowsErr {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").RowError(0, errors.New("rows error"))
				mock.ExpectQuery(query).WillReturnRows(rows)
			} else {
				rows := sqlmock.NewRows([]string{"name"})
				for _, name := range tt.mockNames {
					rows.AddRow(name)
				}
				mock.ExpectQuery(query).WillReturnRows(rows)
			}

			names, err := service.GetNames()
			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expectResult, names)

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBService_GetUniqueNames(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	service := New(db)

	tests := []struct {
		name         string
		mockNames    []string
		mockErr      error
		scanErr      bool
		rowsErr      bool
		expectErr    bool
		expectResult []string
	}{
		{
			name:         "success with duplicate names",
			mockNames:    []string{"Alice", "Bob", "Alice"},
			expectResult: []string{"Alice", "Bob", "Alice"},
		},
		{
			name:         "success with unique names",
			mockNames:    []string{"Alice", "Bob", "Charlie"},
			expectResult: []string{"Alice", "Bob", "Charlie"},
		},
		{
			name:         "success with empty result",
			mockNames:    []string{},
			expectResult: []string{},
		},
		{
			name:      "query error",
			mockErr:   errors.New("query failed"),
			expectErr: true,
		},
		{
			name:         "single name",
			mockNames:    []string{"John"},
			expectResult: []string{"John"},
		},
		{
			name:      "scan error",
			scanErr:   true,
			expectErr: true,
		},
		{
			name:      "rows error",
			rowsErr:   true,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := "SELECT DISTINCT name FROM users"

			if tt.mockErr != nil {
				mock.ExpectQuery(query).WillReturnError(tt.mockErr)
			} else if tt.scanErr {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice")
				rows.RowError(0, errors.New("scan error"))
				mock.ExpectQuery(query).WillReturnRows(rows)
			} else if tt.rowsErr {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").RowError(0, errors.New("rows error"))
				mock.ExpectQuery(query).WillReturnRows(rows)
			} else {
				rows := sqlmock.NewRows([]string{"name"})
				for _, name := range tt.mockNames {
					rows.AddRow(name)
				}
				mock.ExpectQuery(query).WillReturnRows(rows)
			}

			names, err := service.GetUniqueNames()
			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expectResult, names)

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
