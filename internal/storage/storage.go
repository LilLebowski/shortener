package storage

import (
	"github.com/LilLebowski/shortener/internal/storage/db"
	"github.com/LilLebowski/shortener/internal/storage/file"
	"github.com/LilLebowski/shortener/internal/storage/memory"
	"github.com/LilLebowski/shortener/internal/utils"
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
	dbInstance, err := db.Init(databasePath)
	if err != nil {
		utils.Log.Error(err)
	}
	return &Storage{
		File:     file.Init(filePath),
		Memory:   memory.Init(),
		Database: dbInstance,
	}
}
