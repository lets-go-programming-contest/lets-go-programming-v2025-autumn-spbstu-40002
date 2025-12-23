package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"task-6/internal/db"
)

const (
queryGetNames       = "^SELECT name FROM users$"
queryGetUniqueNames = "^SELECT DISTINCT name FROM users$"
)

var errExpected = errors.New("expected error")
	queryGetNames       = "^SELECT name FROM users$"
	queryGetUniqueNames = "^SELECT DISTINCT name FROM users$"
)

var errExpected = errors.New("expected error")

func TestNew(t *testing.T) {
	t.Parallel()

	dbConn, _, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	service := db.New(dbConn)
	require.Equal(t, dbConn, service.DB)
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	type testCase struct {
		rows       *sqlmock.Rows
		queryErr   error
		want       []string
		errIs      error
		errContain string
	}

	cases := []testCase{
		{
			rows: sqlmock.NewRows([]string{"name"}).
				AddRow("Ivan").
				AddRow("Petr"),
			want: []string{"Ivan", "Petr"},
		},
		{
			rows: sqlmock.NewRows([]string{"name"}),
			want: nil,
		},
		{
			queryErr:   errExpected,
			errIs:      errExpected,
			errContain: "db query",
		},
		{
			rows:       sqlmock.NewRows([]string{"name"}).AddRow(nil),
			errContain: "rows scanning",
		},
		{
			rows: sqlmock.NewRows([]string{"name"}).
				AddRow("Ivan").
				AddRow("Petr").
				RowError(1, errExpected),
			errIs:      errExpected,
			errContain: "rows error",
		},
	}

	for i, tc := range cases {
		dbConn, mock, err := sqlmock.New()
		require.NoError(t, err)

		service := db.DBService{DB: dbConn}

		if tc.queryErr != nil {
			mock.ExpectQuery(queryGetNames).WillReturnError(tc.queryErr)
		} else {
			mock.ExpectQuery(queryGetNames).WillReturnRows(tc.rows)
		}

		got, err := service.GetNames()

		if tc.errIs != nil || tc.errContain != "" {
			require.Error(t, err, "case %d", i)

			if tc.errIs != nil {
				require.ErrorIs(t, err, tc.errIs, "case %d", i)
			}

			if tc.errContain != "" {
				require.ErrorContains(t, err, tc.errContain, "case %d", i)
			}

			require.Nil(t, got, "case %d", i)
		} else {
			require.NoError(t, err, "case %d", i)
			require.Equal(t, tc.want, got, "case %d", i)
		}

		mock.ExpectClose()
		require.NoError(t, dbConn.Close())
		require.NoError(t, mock.ExpectationsWereMet(), "case %d", i)
	}
}

func TestGetUniqueNames(t *testing.T) {
	t.Parallel()

	type testCase struct {
		rows       *sqlmock.Rows
		queryErr   error
		want       []string
		errIs      error
		errContain string
	}

	cases := []testCase{
		{
			rows: sqlmock.NewRows([]string{"name"}).
				AddRow("Ivan").
				AddRow("Petr"),
			want: []string{"Ivan", "Petr"},
		},
		{
			rows: sqlmock.NewRows([]string{"name"}),
			want: nil,
		},
		{
			queryErr:   errExpected,
			errIs:      errExpected,
			errContain: "db query",
		},
		{
			rows:       sqlmock.NewRows([]string{"name"}).AddRow(nil),
			errContain: "rows scanning",
		},
		{
			rows: sqlmock.NewRows([]string{"name"}).
				AddRow("Ivan").
				AddRow("Petr").
				RowError(1, errExpected),
			errIs:      errExpected,
			errContain: "rows error",
		},
	}

	for i, tc := range cases {
		dbConn, mock, err := sqlmock.New()
		require.NoError(t, err)

		service := db.DBService{DB: dbConn}

		if tc.queryErr != nil {
			mock.ExpectQuery(queryGetUniqueNames).WillReturnError(tc.queryErr)
		} else {
			mock.ExpectQuery(queryGetUniqueNames).WillReturnRows(tc.rows)
		}

		got, err := service.GetUniqueNames()

		if tc.errIs != nil || tc.errContain != "" {
			require.Error(t, err, "case %d", i)

			if tc.errIs != nil {
				require.ErrorIs(t, err, tc.errIs, "case %d", i)
			}

			if tc.errContain != "" {
				require.ErrorContains(t, err, tc.errContain, "case %d", i)
			}

			require.Nil(t, got, "case %d", i)
		} else {
			require.NoError(t, err, "case %d", i)
			require.Equal(t, tc.want, got, "case %d", i)
		}

		mock.ExpectClose()
		require.NoError(t, dbConn.Close())
		require.NoError(t, mock.ExpectationsWereMet(), "case %d", i)
	}
}
