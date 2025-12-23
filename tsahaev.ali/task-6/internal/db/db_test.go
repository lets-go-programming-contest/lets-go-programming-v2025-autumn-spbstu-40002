package db_test

import (
	"errors"
	"testing"

	"github.com/Tsahaev/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

var (
	errDBDown    = errors.New("db down")
	errIteration = errors.New("iteration error")
	errFatal     = errors.New("fatal error")
)

func TestNew(t *testing.T) {
	t.Parallel()

	service := db.New(nil)
	require.Nil(t, service.DB)
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name          string
		expectedError string
	}{
		{
			name:          "Test without DB",
			expectedError: "database connection is nil",
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			dbService := db.New(nil)

			names, err := dbService.GetNames()

			if test.expectedError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), test.expectedError)
				require.Nil(t, names)
			} else {
				require.NoError(t, err)
				require.NotNil(t, names)
			}
		})
	}
}

func TestGetUniqueNames(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name          string
		expectedError string
	}{
		{
			name:          "Test without DB",
			expectedError: "database connection is nil",
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			dbService := db.New(nil)

			names, err := dbService.GetUniqueNames()

			if test.expectedError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), test.expectedError)
				require.Nil(t, names)
			} else {
				require.NoError(t, err)
				require.NotNil(t, names)
			}
		})
	}
}
