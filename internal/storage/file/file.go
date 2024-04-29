package file

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type URLItem struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Store struct {
	isConfigured bool
	path         string
}

func Init(filePath string) *Store {
	return &Store{
		isConfigured: filePath != "",
		path:         filePath,
	}
}

func (s *Store) Set(full string, short string) error {
	file, err := os.OpenFile(s.path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return fmt.Errorf("storage don't open to write! Error: %s. Path: %s", err, s.path)
	}

	item := URLItem{OriginalURL: full, ShortURL: short}
	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("cannot encode storage item %s", err)
	}
	data = append(data, '\n')
	_, err = file.Write(data)

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	return err
}

func (s *Store) Get(short string) (string, error) {
	file, err := os.OpenFile(s.path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", fmt.Errorf("storage don't open to read! Error: %s. Path: %s", err, s.path)
	}

	r := bufio.NewReader(file)
	line, e := readLine(r)
	var item URLItem
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

func (s *Store) IsConfigured() bool {
	return s.isConfigured
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
