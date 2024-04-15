package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Storage struct {
	URLs map[string]string
}

func NewStorage() *Storage {
	return &Storage{
		URLs: make(map[string]string),
	}
}

func (s *Storage) Set(key string, value string) {
	s.URLs[key] = value
}

func (s *Storage) Get(key string) (string, bool) {
	value, exists := s.URLs[key]
	return value, exists
}

type Memory struct {
	Storage *ShortenerService
	File    *os.File
}
type ShortCollector struct {
	NumberUUID  string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func FillFromStorage(storageInstance *Storage, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	newDecoder := json.NewDecoder(file)
	maxUUID := 0
	for {
		var event ShortCollector
		if err := newDecoder.Decode(&event); err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("error decode JSON:", err)
				break
			}
		}
		maxUUID += 1
		storageInstance.Set(event.OriginalURL, event.ShortURL)
	}
	return nil
}

func WriteFile(storageInstance *Storage, filePath string, BaseURL string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	maxUUID := 0
	for key, value := range storageInstance.URLs {
		shortURL := fmt.Sprintf("%s/%s", BaseURL, key)
		maxUUID += 1
		ShortCollector := ShortCollector{
			strconv.Itoa(maxUUID),
			shortURL,
			value,
		}
		writer := bufio.NewWriter(file)
		err = writeEvent(&ShortCollector, writer)
	}
	return err
}

func writeEvent(ShortCollector *ShortCollector, writer *bufio.Writer) error {
	data, err := json.Marshal(&ShortCollector)
	if err != nil {
		return err
	}

	if _, err := writer.Write(data); err != nil {
		return err
	}

	if err := writer.WriteByte('\n'); err != nil {
		return err
	}

	return writer.Flush()
}
