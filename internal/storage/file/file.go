// Package file contains methods for file storage
package file

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/LilLebowski/shortener/internal/middleware"
	"github.com/LilLebowski/shortener/internal/models"
)

// Storage file storage struct
type Storage struct {
	path string
}

// Init initialization for file storage
func Init(filePath string) *Storage {
	return &Storage{
		path: filePath,
	}
}

// Ping ping storage
func (s *Storage) Ping() error {
	return nil
}

// Set save link info to storage
func (s *Storage) Set(full string, short string, userID string) error {
	file, err := os.OpenFile(s.path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return fmt.Errorf("storage don't open to write! Error: %s. Path: %s", err, s.path)
	}

	item := models.UserURL{OriginalURL: full, ShortURL: short, UserID: userID}
	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("cannot encode storage item %s", err)
	}
	data = append(data, '\n')
	_, err = file.Write(data)

	defer func(file *os.File) {
		errClose := file.Close()
		if errClose != nil {
			middleware.Log.Error("error: close file: %d", errClose)
		}
	}(file)

	return err
}

// SetBatch save batch links info to storage
func (s *Storage) SetBatch(userID string, urls []models.FullURLs) error {
	file, err := os.OpenFile(s.path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return fmt.Errorf("storage don't open to write! Error: %s. Path: %s", err, s.path)
	}

	for _, url := range urls {
		item := models.UserURL{OriginalURL: url.OriginalURL, ShortURL: url.ShortURL, UserID: userID}
		data, errMarshal := json.Marshal(item)
		if errMarshal != nil {
			return fmt.Errorf("cannot encode storage item %s", errMarshal)
		}
		data = append(data, '\n')
		file.Write(data)
	}

	defer func(file *os.File) {
		errClose := file.Close()
		if errClose != nil {
			middleware.Log.Error("error: close file: %d", errClose)
		}
	}(file)

	return err
}

// Get get link info from storage
func (s *Storage) Get(short string) (string, error) {
	file, err := os.OpenFile(s.path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", fmt.Errorf("storage don't open to read! Error: %s. Path: %s", err, s.path)
	}

	r := bufio.NewReader(file)
	line, e := readLine(r)
	var item models.UserURL
	for e == nil {
		err = json.Unmarshal([]byte(line), &item)
		if err != nil {
			return "", fmt.Errorf("storage don't open to read! Error: %s. Path: %s", err, s.path)
		}
		if item.ShortURL == short {
			return item.OriginalURL, nil
		}
		line, e = readLine(r)
	}
	return "", nil
}

// GetByUserID get links info from storage by userID
func (s *Storage) GetByUserID(userID string, baseURL string) ([]map[string]string, error) {
	urls := make([]map[string]string, 0)
	file, err := os.OpenFile(s.path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return urls, fmt.Errorf("storage don't open to read! Error: %s. Path: %s", err, s.path)
	}

	r := bufio.NewReader(file)
	line, e := readLine(r)
	var item models.UserURL
	for e == nil {
		err = json.Unmarshal([]byte(line), &item)
		if err != nil {
			return urls, fmt.Errorf("storage don't open to read! Error: %s. Path: %s", err, s.path)
		}
		if item.UserID == userID {
			shortURL := fmt.Sprintf("%s/%s", baseURL, item.ShortURL)
			urlMap := map[string]string{"short_url": shortURL, "original_url": item.OriginalURL}
			urls = append(urls, urlMap)
		}
		line, e = readLine(r)
	}
	return urls, nil
}

// Delete delete link from storage
func (s *Storage) Delete(string, string, chan<- string) error {
	return nil
}

// Close for closing connection
func (s *Storage) Close() error {
	return nil
}

func readLine(r *bufio.Reader) (string, error) {
	var (
		isPrefix       = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}
