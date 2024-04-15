package utils

import (
	"fmt"

	"github.com/google/uuid"
)

type Repository interface {
	Set(shortID string, originalURL string)
	Get(shortID string) (string, bool)
}

type ShortenerService struct {
	BaseURL string
	Storage Repository
}

func NewShortenerService(BaseURL string, storage Repository) *ShortenerService {
	s := &ShortenerService{
		BaseURL: BaseURL,
		Storage: storage,
	}
	return s
}

func (s *ShortenerService) Set(originalURL string) string {
	shortID := randSeq()
	s.Storage.Set(shortID, originalURL)
	shortURL := fmt.Sprintf("%s/%s", s.BaseURL, shortID)
	return shortURL
}

func randSeq() string {
	newUUID := uuid.New()
	return newUUID.String()
}

func (s *ShortenerService) Get(shortID string) (string, bool) {
	return s.Storage.Get(shortID)
}
