package helpers

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Database struct {
	Conn *sql.DB
}

func New(dbUsername, dbPass, dbPort, dbName, sqlStatement string) (*Database, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		dbUsername, dbPass, dbPort, dbName)

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(); err != nil {
		conn.Close()
		return nil, err
	}

	if _, err = conn.Exec(sqlStatement); err != nil {
		conn.Close()
		return nil, err
	}

	return &Database{Conn: conn}, nil
}
