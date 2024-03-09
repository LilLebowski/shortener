package config

import (
	"flag"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type ShortenerConfiguration struct {
	ServerAddress string
	BaseURL       string
}

func LoadConfiguration() *ShortenerConfiguration {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
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
