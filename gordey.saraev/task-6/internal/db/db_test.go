package db

import (
	"database/sql"
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
			name:         "scan error",
			mockNames:    []string{"Alice"},
			expectErr:    true,
			mockErr:      nil,
			expectResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := "SELECT name FROM users"

			if tt.name == "scan error" {
				rows := sqlmock.NewRows([]string{"name", "extra_column"})
				rows.AddRow("Alice", "extra")
				mock.ExpectQuery(query).WillReturnRows(rows)
			} else if tt.mockErr != nil {
				mock.ExpectQuery(query).WillReturnError(tt.mockErr)
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
			name:         "scan error",
			mockNames:    []string{"Alice"},
			expectErr:    true,
			mockErr:      nil,
			expectResult: nil,
		},
		{
			name:      "rows error",
			mockNames: []string{"Alice", "Bob"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := "SELECT DISTINCT name FROM users"

			if tt.name == "scan error" {
				rows := sqlmock.NewRows([]string{"name", "extra_column"})
				rows.AddRow("Alice", "extra")
				mock.ExpectQuery(query).WillReturnRows(rows)
			} else if tt.name == "rows error" {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").AddRow("Bob").RowError(1, errors.New("row error"))
				mock.ExpectQuery(query).WillReturnRows(rows)
			} else if tt.mockErr != nil {
				mock.ExpectQuery(query).WillReturnError(tt.mockErr)
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

func TestDBService_New(t *testing.T) {
	mockDB := &MockDatabase{}
	service := New(mockDB)
	assert.Equal(t, mockDB, service.DB)
}

type MockDatabase struct {
	QueryFunc func(query string, args ...any) (*sql.Rows, error)
}

func (m *MockDatabase) Query(query string, args ...any) (*sql.Rows, error) {
	if m.QueryFunc != nil {
		return m.QueryFunc(query, args...)
	}
	return nil, nil
}
