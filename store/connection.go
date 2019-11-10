package store

import (
	"database/sql"
	"fmt"

	// Load the sqlite3 package so we can connect.
	_ "github.com/mattn/go-sqlite3"
)

// NewConnection returns a new sqlite3 database connection.
func NewConnection(sqliteConnectString string) (*sql.DB, error) {

	db, err := sql.Open("sqlite3", sqliteConnectString)
	if err != nil {
		return nil, fmt.Errorf("cannot open requested database: %w", err)
	}

	return db, nil
}
