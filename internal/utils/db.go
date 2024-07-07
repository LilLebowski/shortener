// Package utils contains method for db connection
package utils

import (
	"database/sql"
	"fmt"
)

// NewDB creates db connection
func NewDB(databasePath string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databasePath)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}
	return db, nil
}
