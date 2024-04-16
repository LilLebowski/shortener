package storage

import (
	"github.com/LilLebowski/shortener/internal/storage/db"
	"github.com/LilLebowski/shortener/internal/storage/file"
	"github.com/LilLebowski/shortener/internal/storage/memory"
)

type MemoryRepository interface {
	Set(originalURL string, shortID string)
	Get(shortID string) (string, bool)
}

type FileRepository interface {
	Set(full string, short string) error
	Get(short string) (string, error)
}

type DBRepository interface {
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
