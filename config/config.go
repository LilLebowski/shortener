package config

import (
	"flag"
	"github.com/joho/godotenv"
	"os"
	"regexp"
)

const projectDirName = "shortener"

type ShortenerConfiguration struct {
	ServerAddress string
	BaseURL       string
}

func LoadConfiguration() *ShortenerConfiguration {
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		os.Exit(-1)
	}
	cfg := &ShortenerConfiguration{}
	if cfg.ServerAddress = os.Getenv("SERVER_ADDRESS"); cfg.ServerAddress == "" {
		flag.StringVar(&cfg.ServerAddress, "listen", ":8080", "Server address")
	}

	if cfg.BaseURL = os.Getenv("BASE_URL"); cfg.BaseURL == "" {
		flag.StringVar(&cfg.BaseURL, "url", "http://localhost:8080", "Server base URL")
	}
	flag.Parse()

	return cfg
}
