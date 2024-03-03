package utils

import (
	"encoding/base64"
	"fmt"
)

func EncodeUrl(url string) (string, error) {
	fmt.Println("encoded url:", base64.StdEncoding.EncodeToString([]byte(url)))
	shortUrl := base64.StdEncoding.EncodeToString([]byte(url))[:8]
	fmt.Println("shortUrl:", shortUrl)
	return shortUrl, nil
}
