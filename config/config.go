package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type ShortenerConfiguration struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	LogLevel      string `env:"FLAG_LOG_LEVEL"`
	FilePath      string `env:"FILE_STORAGE_PATH"`
	DBPath        string `env:"DATABASE_DSN"`
}

func LoadConfiguration() *ShortenerConfiguration {
	cfg := &ShortenerConfiguration{}
	regStringVar(&cfg.ServerAddress, "a", "localhost:8080", "Server address")
	regStringVar(&cfg.BaseURL, "b", "http://localhost:8080", "Server base URL")
	regStringVar(&cfg.LogLevel, "c", "debug", "Server log level")
	regStringVar(&cfg.FilePath, "f", "/tmp/short-url-db.json", "Server file storage")
	regStringVar(
		&cfg.DBPath,
		"d",
		"host=localhost port=5432 user=admin password=12345 dbname=shortener sslmode=disable",
		"Server db path",
	)
	flag.Parse()

	//fmt.Printf(cfg.FilePath)
	fmt.Printf("base URL: %s\n", cfg.DBPath)

	envServerAddress := os.Getenv("SERVER_ADDRESS")
	envServerAddress = strings.TrimSpace(envServerAddress)
	if envServerAddress != "" {
		cfg.ServerAddress = envServerAddress
	}

	envBaseURL := os.Getenv("BASE_URL")
	envBaseURL = strings.TrimSpace(envBaseURL)
	if envBaseURL != "" {
		cfg.BaseURL = envBaseURL
	}

	envLogLevel := os.Getenv("LOG_LEVEL")
	envLogLevel = strings.TrimSpace(envLogLevel)
	if envLogLevel != "" {
		cfg.LogLevel = envLogLevel
	}

	envFilePath := os.Getenv("FILE_STORAGE_PATH")
	envFilePath = strings.TrimSpace(envFilePath)
	if envFilePath != "" {
		cfg.FilePath = envFilePath
	}

	envDBPath := os.Getenv("DATABASE_DSN")
	envDBPath = strings.TrimSpace(envDBPath)
	if envDBPath != "" {
		cfg.DBPath = envDBPath
	}

	return cfg
}

func regStringVar(p *string, name string, value string, usage string) {
	if flag.Lookup(name) == nil {
		flag.StringVar(p, name, value, usage)
	}
}

func getStringFlag(name string) string {
	return flag.Lookup(name).Value.(flag.Getter).Get().(string)
}
