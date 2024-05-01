package shortener

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/LilLebowski/shortener/internal/storage"
)

type Service struct {
	BaseURL string
	Storage *storage.Storage
}

func Init(BaseURL string, storageInstance *storage.Storage) *Service {
	s := &Service{
		BaseURL: BaseURL,
		Storage: storageInstance,
	}
	return s
}

func (s *Service) Set(originalURL string, userID string) (string, error) {
	shortID := getShortURL(originalURL)
	shortURL := fmt.Sprintf("%s/%s", s.BaseURL, shortID)
	if s.Storage.Database.IsConfigured() {
		err := s.Storage.Database.Set(originalURL, shortID, userID)
		if err != nil {
			return shortURL, err
		}
	} else if s.Storage.File.IsConfigured() {
		err := s.Storage.File.Set(originalURL, shortID, userID)
		if err != nil {
			return shortURL, err
		}
	}
	s.Storage.Memory.Set(originalURL, shortID, userID)
	return shortURL, nil
}

func (s *Service) Get(shortID string) (string, bool) {
	if s.Storage.Database.IsConfigured() {
		fullURL, _ := s.Storage.Database.Get(shortID)
		if fullURL != "" {
			return fullURL, true
		}
	} else if s.Storage.File.IsConfigured() {
		fullURL, _ := s.Storage.File.Get(shortID)
		if fullURL != "" {
			return fullURL, true
		}
	}
	return s.Storage.Memory.Get(shortID)
}

func (s *Service) Ping() error {
	return s.Storage.Database.Ping()
}

func (s *Service) GetByUserID(userID string) ([]map[string]string, error) {
	if s.Storage.Database.IsConfigured() {
		urls, err := s.Storage.Database.GetByUserID(userID, s.BaseURL)
		return urls, err
	} else if s.Storage.File.IsConfigured() {
		urls, err := s.Storage.File.GetByUserID(userID, s.BaseURL)
		return urls, err
	}
	return s.Storage.Memory.GetByUserID(userID, s.BaseURL)
}

func getShortURL(longURL string) string {
	splitURL := strings.Split(longURL, "://")
	hash := sha1.New()
	if len(splitURL) < 2 {
		hash.Write([]byte(longURL))
	} else {
		hash.Write([]byte(splitURL[1]))
	}
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}
