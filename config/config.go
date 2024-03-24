package config

import (
	"flag"
	"os"
)

type ShortenerConfiguration struct {
	ServerAddress string
	BaseURL       string
}

func LoadConfiguration() *ShortenerConfiguration {
	cfg := &ShortenerConfiguration{}
	if cfg.ServerAddress = os.Getenv("SERVER_ADDRESS"); cfg.ServerAddress == "" {
		flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "Server address")
	}

	if cfg.BaseURL = os.Getenv("BASE_URL"); cfg.BaseURL == "" {
		flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "Server base URL")
	}
	flag.Parse()

	return cfg
}
