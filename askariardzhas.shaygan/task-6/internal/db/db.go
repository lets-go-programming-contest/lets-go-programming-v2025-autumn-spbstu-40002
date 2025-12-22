package db

import (
	"database/sql"
	"fmt"
)

type DatabaseConn interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type DBClient struct {
	Conn DatabaseConn
}

func Create(conn DatabaseConn) DBClient {
	return DBClient{Conn: conn}
}

func (client DBClient) GetAllNames() ([]string, error) {
	sql := "SELECT name FROM users"

	rows, err := client.Conn.Query(sql)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	var result []string

	for rows.Next() {
		var n string

		if err := rows.Scan(&n); err != nil {
			return nil, fmt.Errorf("scanning row failed: %w", err)
		}

		result = append(result, n)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return result, nil
}

func (client DBClient) GetDistinctNames() ([]string, error) {
	sql := "SELECT DISTINCT name FROM users"

	rows, err := client.Conn.Query(sql)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	var uniqueNames []string

	for rows.Next() {
		var name string

		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("scanning row failed: %w", err)
		}

		uniqueNames = append(uniqueNames, name)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return uniqueNames, nil
}
