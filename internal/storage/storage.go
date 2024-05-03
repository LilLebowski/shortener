package storage

import (
	"github.com/LilLebowski/shortener/internal/storage/db"
	"github.com/LilLebowski/shortener/internal/storage/file"
	"github.com/LilLebowski/shortener/internal/storage/memory"
)

type MemoryRepository interface {
	Set(full string, short string)
	Get(shortID string) (string, bool)
}

type FileRepository interface {
	Set(full string, short string) error
	Get(short string) (string, error)
	IsConfigured() bool
}

type DBRepository interface {
	Set(full string, short string) error
	Get(short string) (string, error)
	IsConfigured() bool
	Ping() error
}

type Storage struct {
	File     FileRepository
	Database DBRepository
	Memory   MemoryRepository
}

func Init(filePath string, databasePath string) *Storage {
	// FIXME: вынести в service/shortener и избавиться от IsConfigured и лишних интерфейсов
	dbInstance, _ := db.Init(databasePath)
	return &Storage{
		File:     file.Init(filePath),
		Memory:   memory.Init(),
		Database: dbInstance,
	}
}
