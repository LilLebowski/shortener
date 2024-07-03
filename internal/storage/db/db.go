package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/LilLebowski/shortener/internal/middleware"
	"github.com/LilLebowski/shortener/internal/models"
	"go.uber.org/zap"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/LilLebowski/shortener/internal/utils"
)

type Storage struct {
	db *sql.DB
}

func Init(db *sql.DB) (*Storage, error) {
	err := createTable(db)
	if err != nil {
		return nil, fmt.Errorf("error create table db: %w", err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) Set(full string, short string, userID string) error {
	query := `
        INSERT INTO url (short_id, original_url, user_id) 
        VALUES ($1, $2, $3)
    `
	_, err := s.db.Exec(query, short, full, userID)
	var e *pgconn.PgError
	if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
		return utils.NewUniqueConstraintError(err)
	}
	return err
}

func (s *Storage) SetBatch(userID string, urls []models.FullURLs) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	query := `
		INSERT INTO url (short_id, original_url, user_id) 
		VALUES (DEFAULT, $1, $2, $3);
	`

	stmt, err := tx.PrepareContext(context.Background(), query)

	defer func(stmt *sql.Stmt) {
		err = stmt.Close()
		if err != nil {
			middleware.Log.Info("Close statement error", zap.Error(err))
		}
	}(stmt)

	for _, v := range urls {
		_, err = stmt.ExecContext(context.Background(), v.ShortURL, v.OriginalURL, userID)
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			return utils.NewUniqueConstraintError(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Get(short string) (string, error) {
	query := `
        SELECT original_url, is_deleted
        FROM url 
        WHERE short_id = $1
    `
	var originalURL string
	var isDeleted bool
	err := s.db.QueryRow(query, short).Scan(&originalURL, &isDeleted)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", utils.NewNotFoundError("url not found", err)
		}
		return "", err
	}

	if isDeleted {
		return "", utils.NewDeletedError("url is already deleted", err)
	}

	return originalURL, nil
}

func (s *Storage) GetByUserID(userID string, baseURL string) ([]map[string]string, error) {
	urls := make([]map[string]string, 0)
	query := `SELECT original_url, short_id FROM url WHERE user_id=$1 AND is_deleted=false;`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return urls, err
	}
	if rows.Err() != nil {
		return urls, rows.Err()
	}
	defer rows.Close()

	for rows.Next() {
		var shortID, originalURL string
		if err = rows.Scan(&originalURL, &shortID); err != nil {
			return nil, err
		}
		shortURL := fmt.Sprintf("%s/%s", baseURL, shortID)
		urlMap := map[string]string{"short_url": shortURL, "original_url": originalURL}
		urls = append(urls, urlMap)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during iteration through link rows: %w", err)
	}

	return urls, nil
}

func (s *Storage) Delete(userID string, shortURL string, updateChan chan<- string) error {
	query := `
		UPDATE url
		SET is_deleted = true
		WHERE short_id = $1 and  user_id = $2`

	_, err := s.db.Exec(query, shortURL, userID)
	if err != nil {
		return err
	}
	updateChan <- shortURL
	return nil
}

func (s *Storage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("pinging db-store: %w", err)
	}
	return nil
}

func createTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS url (
		id SERIAL PRIMARY KEY,
		short_id VARCHAR(256) NOT NULL UNIQUE,
		original_url TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    	user_id VARCHAR(360),
		is_deleted BOOLEAN NOT NULL DEFAULT FALSE
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
