package memory

import (
	"fmt"

	"github.com/LilLebowski/shortener/internal/models"
	"github.com/LilLebowski/shortener/internal/utils"
)

type URLItem struct {
	OriginalURL string
	UserID      string
}

type Storage struct {
	URLs map[string]*URLItem
}

func Init() *Storage {
	return &Storage{
		URLs: make(map[string]*URLItem),
	}
}

func (s *Storage) Ping() error {
	return nil
}

func (s *Storage) Set(full string, short string, userID string) error {
	s.URLs[short] = &URLItem{
		OriginalURL: full,
		UserID:      userID,
	}
	return nil
}

func (s *Storage) SetBatch(userID string, urls []models.FullURLs) error {
	for _, url := range urls {
		s.URLs[url.ShortURL] = &URLItem{
			OriginalURL: url.OriginalURL,
			UserID:      userID,
		}
	}
	return nil
}

func (s *Storage) Get(short string) (string, error) {
	value, exists := s.URLs[short]
	if exists {
		return value.OriginalURL, nil
	}
	return "", utils.NewDeletedError("url is already deleted", nil)
}

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

func (s *Storage) Delete(userID string, shortURL string, updateChan chan<- string) error {
	return nil
}
