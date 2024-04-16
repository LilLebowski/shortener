package config

import (
	"flag"
	"github.com/caarlos0/env"
	"log"
)

const (
	FileName      = "/tmp/short-url-db.json"
	ServerAddress = "localhost:8080"
	BaseURL       = "http://localhost:8080"
	DBPath        = "host=localhost port=5432 user=admin password=12345 dbname=shortener sslmode=disable"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	LogLevel      string `env:"FLAG_LOG_LEVEL"`
	FilePath      string `env:"FILE_STORAGE_PATH"`
	DBPath        string `env:"DATABASE_DSN"`
}

func LoadConfiguration() *Config {
	flagServerAddress := flag.String("a", ServerAddress, "server adress")
	flagBaseURL := flag.String("b", BaseURL, "base url")
	flagFilePath := flag.String("c", FileName, "file path")
	flagDataBaseURI := flag.String("d", DBPath, "URI for database")
	flag.Parse()

	cfg := Config{
		ServerAddress: ServerAddress,
		FilePath:      FileName,
		BaseURL:       BaseURL,
		DBPath:        DBPath,
	}

	err := env.Parse(&cfg)

	if err != nil {
		log.Fatal(err)
	}

	if *flagServerAddress != ServerAddress {
		cfg.ServerAddress = *flagServerAddress
	}
	if *flagBaseURL != BaseURL {
		cfg.BaseURL = *flagBaseURL
	}
	if *flagFilePath != FileName {
		cfg.FilePath = *flagFilePath
	}
	if *flagDataBaseURI != DBPath {
		cfg.DBPath = *flagDataBaseURI
	}

	return &cfg
}
