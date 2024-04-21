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

func (s *Service) Set(originalURL string) (string, error) {
	shortID := getShortURL(originalURL)
	shortURL := fmt.Sprintf("%s/%s", s.BaseURL, shortID)
	if s.Storage.Database.IsConfigured() {
		err := s.Storage.Database.Set(originalURL, shortID)
		if err != nil {
			return shortURL, err
		}
	} else if s.Storage.File.IsConfigured() {
		err := s.Storage.File.Set(originalURL, shortID)
		if err != nil {
			return shortURL, err
		}
	}
	s.Storage.Memory.Set(originalURL, shortID)
	return shortURL, nil
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
