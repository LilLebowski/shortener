package storage

import (
	"github.com/LilLebowski/shortener/internal/storage/db"
	"github.com/LilLebowski/shortener/internal/storage/file"
	"github.com/LilLebowski/shortener/internal/storage/memory"
)

type MemoryRepository interface {
	Set(full string, short string, userID string)
	Get(shortID string) (string, bool)
	GetByUserID(userID string, baseURL string) ([]map[string]string, error)
}

type FileRepository interface {
	Set(full string, short string, userID string) error
	Get(short string) (string, error)
	GetByUserID(userID string, baseURL string) ([]map[string]string, error)
	IsConfigured() bool
}

type DBRepository interface {
	Set(full string, short string, userID string) error
	Get(short string) (string, bool, error)
	GetByUserID(userID string, baseURL string) ([]map[string]string, error)
	Delete(userID string, shortURL string, updateChan chan<- string) error
	IsConfigured() bool
	Ping() error
}

type Storage struct {
	File     FileRepository
	Database DBRepository
	Memory   MemoryRepository
}

func Init(filePath string, databasePath string) *Storage {
	dbInstance, _ := db.Init(databasePath)
	return &Storage{
		File:     file.Init(filePath),
		Memory:   memory.Init(),
		Database: dbInstance,
	}
}
