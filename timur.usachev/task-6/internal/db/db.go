package db

import "database/sql"

type Store struct {
	DB *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{DB: db}
}

var closeRows = func(r *sql.Rows) error {
	return r.Close()
}

func (s *Store) GetNames() (res []string, err error) {
	rows, err := s.DB.Query("SELECT name FROM users")
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := closeRows(rows); cerr != nil && err == nil {
			err = cerr
		}
	}()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		res = append(res, name)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}
