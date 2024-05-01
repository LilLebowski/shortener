package memory

import "fmt"

type URLItem struct {
	OriginalURL string
	UserID      string
}

type Storage struct {
	URLs map[string]*URLItem
}

func Init() *Storage {
	return &Storage{
		URLs: make(map[string]*URLItem),
	}
}

func (s *Storage) Set(full string, short string, userID string) {
	s.URLs[short] = &URLItem{
		OriginalURL: full,
		UserID:      userID,
	}
}

func (s *Storage) Get(short string) (string, bool) {
	value, exists := s.URLs[short]
	if exists {
		return value.OriginalURL, exists
	}
	return "", false
}

func (s *Storage) GetByUserID(userID string, baseURL string) ([]map[string]string, error) {
	urls := make([]map[string]string, 0)
	for shortID, item := range s.URLs {
		if item.UserID == userID {
			shortURL := fmt.Sprintf("%s/%s", baseURL, shortID)
			urlMap := map[string]string{"short_url": shortURL, "original_url": item.OriginalURL}
			urls = append(urls, urlMap)
		}
	}
	return urls, nil
}
