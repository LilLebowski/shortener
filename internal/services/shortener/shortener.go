package shortener

import (
	"fmt"
	"github.com/LilLebowski/shortener/internal/storage"
	"github.com/google/uuid"
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

func (s *Service) Set(originalURL string) string {
	shortID := randSeq()
	if s.Storage.File.IsConfigured() {
		err := s.Storage.File.Set(originalURL, shortID)
		if err != nil {
			return ""
		}
	}
	s.Storage.Memory.Set(originalURL, shortID)
	shortURL := fmt.Sprintf("%s/%s", s.BaseURL, shortID)
	return shortURL
}

func randSeq() string {
	newUUID := uuid.New()
	return newUUID.String()
}

func (s *Service) Get(shortID string) (string, bool) {
	if s.Storage.File.IsConfigured() {
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
