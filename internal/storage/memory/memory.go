// Package memory contains methods for memory storage
package memory

import (
	"fmt"

	"github.com/LilLebowski/shortener/internal/models"
	"github.com/LilLebowski/shortener/internal/utils"
)

// URLItem memory storage link info
type URLItem struct {
	OriginalURL string
	UserID      string
}

// Storage memory storage struct
type Storage struct {
	URLs map[string]*URLItem
}

// Init initialization for memory storage
func Init() *Storage {
	return &Storage{
		URLs: make(map[string]*URLItem),
	}
}

// Ping ping storage
func (s *Storage) Ping() error {
	return nil
}

// Set save link info to storage
func (s *Storage) Set(full string, short string, userID string) error {
	s.URLs[short] = &URLItem{
		OriginalURL: full,
		UserID:      userID,
	}
	return nil
}

// SetBatch save batch links info to storage
func (s *Storage) SetBatch(userID string, urls []models.FullURLs) error {
	for _, url := range urls {
		s.URLs[url.ShortURL] = &URLItem{
			OriginalURL: url.OriginalURL,
			UserID:      userID,
		}
	}
	return nil
}

// Get get link info from storage
func (s *Storage) Get(short string) (string, error) {
	value, exists := s.URLs[short]
	if exists {
		return value.OriginalURL, nil
	}
	return "", utils.NewDeletedError("url is already deleted", nil)
}

// GetByUserID get links info from storage by userID
func (s *Storage) GetByUserID(userID string, baseURL string) ([]map[string]string, error) {
	urls := make([]map[string]string, 0)
	for shortID, item := range s.URLs {
		if item.UserID == userID {
			shortURL := fmt.Sprintf("%s/%s", baseURL, shortID)
			urlMap := map[string]string{"short_url": shortURL, "original_url": item.OriginalURL}
			urls = append(urls, urlMap)
		}
	}
	return urls, nil
}

// Delete delete link from storage
func (s *Storage) Delete(userID string, shortURL string, updateChan chan<- string) error {
	return nil
}

// Close for closing connection
func (s *Storage) Close() error {
	return nil
}
