package shortener

import (
	"crypto/sha1"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/LilLebowski/shortener/internal/storage"
	"github.com/LilLebowski/shortener/internal/utils"
	"go.uber.org/zap"
	"log"
	"strings"

	"github.com/LilLebowski/shortener/cmd/shortener/config"
	"github.com/LilLebowski/shortener/internal/middleware"
	dbs "github.com/LilLebowski/shortener/internal/storage/db"
	fs "github.com/LilLebowski/shortener/internal/storage/file"
	ms "github.com/LilLebowski/shortener/internal/storage/memory"
)

type Service struct {
	BaseURL string
	Storage storage.Repository
}

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
	} else if config.FilePath != "" {
		s.Storage = fs.Init(config.FilePath)
	} else {
		s.Storage = ms.Init()
	}
	return s
}

func (s *Service) Set(originalURL string, userID string) (string, error) {
	shortID := getShortURL(originalURL)
	shortURL := fmt.Sprintf("%s/%s", s.BaseURL, shortID)
	err := s.Storage.Set(originalURL, shortID, userID)
	if err != nil {
		return shortURL, err
	}
	return shortURL, nil
}

func (s *Service) Get(shortID string) (string, bool, bool) {
	var deletedErr *utils.DeletedError
	var notFoundErr *utils.NotFoundError
	fullURL, err := s.Storage.Get(shortID)
	if errors.As(err, &deletedErr) {
		return fullURL, false, false
	}
	if errors.As(err, &notFoundErr) {
		return fullURL, true, false
	}
	return fullURL, false, true
}

func (s *Service) Ping() error {
	return s.Storage.Ping()
}

func (s *Service) GetByUserID(userID string) ([]map[string]string, error) {
	return s.Storage.GetByUserID(userID, s.BaseURL)
}

func (s *Service) DeleteURLsRep(userID string, shorURLs []string) error {
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

func getShortURL(longURL string) string {
	splitURL := strings.Split(longURL, "://")
	hash := sha1.New()
	if len(splitURL) < 2 {
		hash.Write([]byte(longURL))
	} else {
		hash.Write([]byte(splitURL[1]))
	}
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}
