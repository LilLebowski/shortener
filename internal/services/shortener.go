package services

import (
	"fmt"
	"github.com/LilLebowski/shortener/internal/storage"
	"github.com/google/uuid"
)

type ShortenerService struct {
	BaseURL string
	Storage *storage.Storage
}

func Init(BaseURL string, storageInstance *storage.Storage) *ShortenerService {
	s := &ShortenerService{
		BaseURL: BaseURL,
		Storage: storageInstance,
	}
	return s
}

func (s *ShortenerService) Set(originalURL string) string {
	shortID := randSeq()
	if s.Storage.File != nil {
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

func (s *ShortenerService) Get(shortID string) (string, bool) {
	if s.Storage.File != nil {
		fullURL, _ := s.Storage.File.Get(shortID)
		if fullURL != "" {
			return fullURL, true
		}
	}
	return s.Storage.Memory.Get(shortID)
}

func (s *ShortenerService) Ping() error {
	return s.Storage.Database.Ping()
}
