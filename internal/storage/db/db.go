package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/LilLebowski/shortener/internal/utils"
)

type Store struct {
	isConfigured bool
	db           *sql.DB
}

func Init(databasePath string) (*Store, error) {
	db, err := sql.Open("pgx", databasePath)
	dbStore := &Store{
		isConfigured: databasePath != "",
	}
	if err != nil {
		dbStore.isConfigured = false
		return dbStore, fmt.Errorf("error opening db: %w", err)
	}

	err = createTable(db)
	if err != nil {
		dbStore.isConfigured = false
		fmt.Printf("error create table db: %w", err)
		return dbStore, fmt.Errorf("error create table db: %w", err)
	}

	dbStore.db = db

	return dbStore, nil
}

func (s *Store) Set(full string, short string) error {
	query := `
        INSERT INTO url (short_id, original_url) 
        VALUES ($1, $2)
    `
	_, err := s.db.Exec(query, short, full)
	var e *pgconn.PgError
	if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
		return utils.NewUniqueConstraintError(err)
	}
	return err
}

func (s *Store) Get(short string) (string, error) {
	query := `
        SELECT original_url 
        FROM url 
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
	);
	DO $$ 
		BEGIN 
		 IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE tablename = 'url' AND indexname = 'idx_original_url') THEN
			CREATE UNIQUE INDEX idx_original_url ON url(original_url);
		END IF;
	END $$;
`

	_, err := db.Exec(query)
	return err
}
