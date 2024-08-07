// Package shortener contains methods for shortener service
package shortener

import (
	"crypto/sha1"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"log"
	"strings"

	"github.com/LilLebowski/shortener/cmd/shortener/config"
	"github.com/LilLebowski/shortener/internal/middleware"
	"github.com/LilLebowski/shortener/internal/models"
	"github.com/LilLebowski/shortener/internal/storage"
	dbs "github.com/LilLebowski/shortener/internal/storage/db"
	fs "github.com/LilLebowski/shortener/internal/storage/file"
	ms "github.com/LilLebowski/shortener/internal/storage/memory"
	"github.com/LilLebowski/shortener/internal/utils"
)

// Service shortener service struct
type Service struct {
	BaseURL string
	Storage storage.Repository
}

// Init initialization for shortener service
func Init(config *config.Config) *Service {
	s := &Service{
		BaseURL: config.BaseURL,
	}
	var dbc *sql.DB
	var err error
	if config.DBPath != "" {
		dbc, err = utils.NewDB(config.DBPath)
	}
	if config.DBPath != "" && err == nil {
		s.Storage, err = dbs.Init(dbc)
		if err != nil {
			log.Fatal(err)
		}
		return s
	}
	if config.FilePath != "" {
		s.Storage = fs.Init(config.FilePath)
	} else {
		s.Storage = ms.Init()
	}
	return s
}

// Set create short link and save info to storage
func (s *Service) Set(originalURL string, userID string) (string, error) {
	shortID := GetShortURL(originalURL)
	shortURL := fmt.Sprintf("%s/%s", s.BaseURL, shortID)
	err := s.Storage.Set(originalURL, shortID, userID)
	if err != nil {
		return shortURL, err
	}
	return shortURL, nil
}

// SetBatch create short link and save info to storage from array
func (s *Service) SetBatch(urls []models.URLs, userID string) ([]models.ShortURLs, error) {
	var fullURLs []models.FullURLs
	var shorts []models.ShortURLs

	for _, req := range urls {
		url := strings.TrimSpace(req.OriginalURL)
		shortID := GetShortURL(url)
		shortURL := fmt.Sprintf("%s/%s", s.BaseURL, shortID)
		fullURLs = append(fullURLs, models.FullURLs{
			ShortURL:    shortID,
			OriginalURL: url,
		})
		shorts = append(shorts, models.ShortURLs{
			ShortURL:      shortURL,
			CorrelationID: req.CorrelationID,
		})
	}

	err := s.Storage.SetBatch(userID, fullURLs)

	if err != nil {
		return shorts, err
	}
	return shorts, nil
}

// Get get link info from storage
func (s *Service) Get(shortID string) (string, bool, bool) {
	var deletedErr *utils.DeletedError
	var notFoundErr *utils.NotFoundError
	fullURL, err := s.Storage.Get(shortID)
	if errors.As(err, &deletedErr) {
		return fullURL, true, true
	}
	if errors.As(err, &notFoundErr) {
		return fullURL, false, false
	}
	return fullURL, false, true
}

// Ping ping storage
func (s *Service) Ping() error {
	return s.Storage.Ping()
}

// GetByUserID get links info from storage by userID
func (s *Service) GetByUserID(userID string) ([]map[string]string, error) {
	return s.Storage.GetByUserID(userID, s.BaseURL)
}

// DeleteURLs delete links from storage
func (s *Service) DeleteURLs(userID string, shorURLs []string) error {
	resultChan := make(chan string)
	updateChan := make(chan string, len(shorURLs))

	go func() {
		for _, shortURL := range shorURLs {
			err := s.Storage.Delete(userID, shortURL, updateChan)
			if err != nil {
				middleware.Log.Error("Failed to delete URLs", zap.Error(err))
			}
		}
	}()

	go func() {
		for updateShortID := range updateChan {
			resultChan <- updateShortID
		}
		close(resultChan)
	}()
	return nil
}

// GetShortURL get short link from full
func GetShortURL(fullURL string) string {
	splitURL := strings.Split(fullURL, "://")
	hash := sha1.New()
	if len(splitURL) < 2 {
		hash.Write([]byte(fullURL))
	} else {
		hash.Write([]byte(splitURL[1]))
	}
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}
