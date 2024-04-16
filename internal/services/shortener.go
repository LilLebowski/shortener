package services

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/LilLebowski/shortener/internal/storage"
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
	} else {
		s.Storage.Memory.Set(originalURL, shortID)
	}
	shortURL := fmt.Sprintf("%s/%s", s.BaseURL, shortID)
	return shortURL
}

func randSeq() string {
	newUUID := uuid.New()
	return newUUID.String()
}

func (s *ShortenerService) Get(shortID string) (string, bool) {
	if s.Storage.File != nil {
		fmt.Printf("file!!!")
		fullURL, err := s.Storage.File.Get(shortID)
		if err != nil {
			return fullURL, true
		}
	} else {
		return s.Storage.Memory.Get(shortID)
	}
	return "", false
}

func (s *ShortenerService) Ping() error {
	return s.Storage.Database.Ping()
}
