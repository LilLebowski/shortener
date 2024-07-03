package storage

import "github.com/LilLebowski/shortener/internal/models"

type Repository interface {
	Ping() error
	Set(full string, short string, userID string) error
	SetBatch(userID string, urls []models.FullURLs) error
	Get(short string) (string, error)
	GetByUserID(userID string, baseURL string) ([]map[string]string, error)
	Delete(userID string, shortURL string, updateChan chan<- string) error
	//DeleteBatch(userID string, urls []string) error
}
