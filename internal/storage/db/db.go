package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type Store struct {
	isConfigured bool
	db           *sql.DB
}

func Init(databasePath string) (*Store, error) {
	db, err := sql.Open("pgx", databasePath)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	err = createTable(db)
	if err != nil {
		return nil, fmt.Errorf("error creae table db: %w", err)
	}

	return &Store{
		isConfigured: databasePath != "",
		db:           db,
	}, nil
}

func (s *Store) Set(full string, short string) error {
	query := `
        INSERT INTO urls (short_id, original_url) 
        VALUES ($1, $2)
    `
	_, err := s.db.Exec(query, short, full)
	if err != nil {
		return fmt.Errorf("error save URL: %w", err)
	}
	return nil
}

func (s *Store) Get(short string) (string, error) {
	query := `
        SELECT original_url 
        FROM urls 
        WHERE short_id = $1
    `

	var originalURL string
	err := s.db.QueryRow(query, short).Scan(&originalURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", err
		}
		return "", err
	}

	return originalURL, err
}

func (s *Store) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("pinging db-store: %w", err)
	}
	return nil
}

func (s *Store) IsConfigured() bool {
	return s.isConfigured
}

func createTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS url (
		id SERIAL PRIMARY KEY,
		short_id VARCHAR(256) NOT NULL UNIQUE,
		original_url TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
