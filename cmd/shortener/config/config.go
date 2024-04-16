package config

import (
	"flag"
	"github.com/caarlos0/env"
	"log"
)

const (
	ServerAddress = "localhost:8080"
	BaseURL       = "http://localhost:8080"
	LogLevel      = "debug"
	FileName      = "/tmp/short-url-db.json"
	DBPath        = ""
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	LogLevel      string `env:"FLAG_LOG_LEVEL"`
	FilePath      string `env:"FILE_STORAGE_PATH"`
	DBPath        string `env:"DATABASE_DSN"`
}

func LoadConfiguration() *Config {
	cfg := Config{}

	regStringVar(&cfg.ServerAddress, "a", ServerAddress, "Server address")
	regStringVar(&cfg.BaseURL, "b", BaseURL, "Server base URL")
	regStringVar(&cfg.LogLevel, "c", LogLevel, "Server log level")
	regStringVar(&cfg.FilePath, "f", FileName, "Server file storage")
	regStringVar(&cfg.DBPath, "d", DBPath, "Server db path")

	flag.Parse()

	flagServerAddress := getStringFlag("a")
	flagBaseURL := getStringFlag("b")
	flagFilePath := getStringFlag("f")
	flagDataBaseURI := getStringFlag("d")

	err := env.Parse(&cfg)

	if err != nil {
		log.Fatal(err)
	}

	if flagServerAddress != ServerAddress {
		cfg.ServerAddress = flagServerAddress
	}
	if flagBaseURL != BaseURL {
		cfg.BaseURL = flagBaseURL
	}
	if flagFilePath != FileName {
		cfg.FilePath = flagFilePath
	}
	if flagDataBaseURI != DBPath {
		cfg.DBPath = flagDataBaseURI
	}

	return &cfg
}

func regStringVar(p *string, name string, value string, usage string) {
	if flag.Lookup(name) == nil {
		flag.StringVar(p, name, value, usage)
	}
}

func getStringFlag(name string) string {
	return flag.Lookup(name).Value.(flag.Getter).Get().(string)
}
