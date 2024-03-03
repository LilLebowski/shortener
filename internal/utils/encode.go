package utils

import (
	"encoding/base64"
	"fmt"
)

func EncodeURL(url string) (string, error) {
	fmt.Println("encoded url:", base64.StdEncoding.EncodeToString([]byte(url)))
	shortURL := base64.StdEncoding.EncodeToString([]byte(url))[:8]
	fmt.Println("shortURL:", shortURL)
	return shortURL, nil
}
