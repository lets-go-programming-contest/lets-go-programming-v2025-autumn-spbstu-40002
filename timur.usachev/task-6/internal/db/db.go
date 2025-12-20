package db

import (
	"fmt"
)

type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
	Close() error
}

type Database interface {
	Query(query string, args ...any) (Rows, error)
}

type Store struct {
	DB Database
}

func New(db Database) *Store {
	return &Store{DB: db}
}

func (s *Store) GetNames() (names []string, err error) {
	names = []string{}
	rows, err := s.DB.Query("SELECT name FROM users")
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	defer func() {
		if cerr := rows.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("close: %w", cerr)
		}
	}()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		names = append(names, name)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return names, nil
}

func (s *Store) GetUniqueNames() (names []string, err error) {
	names = []string{}
	rows, err := s.DB.Query("SELECT DISTINCT name FROM users")
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	defer func() {
		if cerr := rows.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("close: %w", cerr)
		}
	}()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		names = append(names, name)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return names, nil
}
