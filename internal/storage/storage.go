// Package storage implement interface for Repository
package storage

import "github.com/LilLebowski/shortener/internal/models"

// Repository  interface for working with global repository
type Repository interface {
	// Ping ping storage
	Ping() error
	// Set save link info to storage
	Set(full string, short string, userID string) error
	// SetBatch save batch links info to storage
	SetBatch(userID string, urls []models.FullURLs) error
	// Get get link info from storage
	Get(short string) (string, error)
	// GetByUserID get links info from storage by userID
	GetByUserID(userID string, baseURL string) ([]map[string]string, error)
	// Delete delete link from storage
	Delete(userID string, shortURL string, updateChan chan<- string) error
}
