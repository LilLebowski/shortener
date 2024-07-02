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
	UserID      string `json:"user_id"`
}

type Storage struct {
	path string
}

func Init(filePath string) *Storage {
	return &Storage{
		path: filePath,
	}
}

func (s *Storage) Ping() error {
	return nil
}

func (s *Storage) Set(full string, short string, userID string) error {
	file, err := os.OpenFile(s.path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return fmt.Errorf("storage don't open to write! Error: %s. Path: %s", err, s.path)
	}

	item := URLItem{OriginalURL: full, ShortURL: short, UserID: userID}
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

func (s *Storage) Get(short string) (string, error) {
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

func (s *Storage) GetByUserID(userID string, baseURL string) ([]map[string]string, error) {
	urls := make([]map[string]string, 0)
	file, err := os.OpenFile(s.path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return urls, fmt.Errorf("storage don't open to read! Error: %s. Path: %s", err, s.path)
	}

	r := bufio.NewReader(file)
	line, e := readLine(r)
	var item URLItem
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

func (s *Storage) Delete(string, string, chan<- string) error {
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
