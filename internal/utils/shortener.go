package utils

import (
	"fmt"

	"github.com/google/uuid"
)

type Repository interface {
	Set(shortID string, originalURL string)
	Get(shortID string) (string, bool)
}

type Store interface {
	PingDBStore() error
}

type ShortenerService struct {
	BaseURL string
	Storage Repository
	DB      Store
}

func NewShortenerService(BaseURL string, storage Repository, db Store) *ShortenerService {
	s := &ShortenerService{
		BaseURL: BaseURL,
		Storage: storage,
		DB:      db,
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

func (s *ShortenerService) Ping() error {
	return s.DB.PingDBStore()
}
