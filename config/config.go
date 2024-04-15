package config

import (
	"flag"
	"os"
	"strings"
)

type ShortenerConfiguration struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	LogLevel      string `env:"FLAG_LOG_LEVEL"`
	FilePath      string `env:"FILE_STORAGE_PATH"`
}

func LoadConfiguration() *ShortenerConfiguration {
	cfg := &ShortenerConfiguration{}
	regStringVar(&cfg.ServerAddress, "a", "localhost:8080", "Server address")
	regStringVar(&cfg.BaseURL, "b", "http://localhost:8080", "Server base URL")
	regStringVar(&cfg.LogLevel, "c", "debug", "Server log level")
	regStringVar(&cfg.FilePath, "f", "short-url-db.json", "Server file storage")
	flag.Parse()

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

	envFilePath := os.Getenv("FILE_PATH")
	envFilePath = strings.TrimSpace(envFilePath)
	if envFilePath != "" {
		cfg.FilePath = envFilePath
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
