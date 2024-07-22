// Package config implement functions for environment and project configs
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/caarlos0/env"
)

// default config constants
const (
	ServerAddress = "localhost:8080"
	BaseURL       = "http://localhost:8080"
	LogLevel      = "debug"
	FileName      = "/tmp/short-url-db.json"
	DBPath        = ""
	TokenExpire   = time.Hour * 24
	SecretKey     = "09d25e094faa6ca2556c818166b7a9563b93f7099f6f0f4caa6cf63b88e8d3e7"
	EnableHTTPS   = ""
)

// Config struct for environment
type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	LogLevel      string `env:"FLAG_LOG_LEVEL"`
	FilePath      string `env:"FILE_STORAGE_PATH"`
	DBPath        string `env:"DATABASE_DSN"`
	EnableHTTPS   string `env:"ENABLE_HTTPS" envDefault:""`

	TokenExpire time.Duration
	SecretKey   string
}

// JSONConfig for json config
type JSONConfig struct {
	BaseURL       string `json:"base_url"`
	ServerAddress string `json:"server_address"`
	FilePath      string `json:"file_storage_path"`
	DBPath        string `json:"database_dsn"`
	EnableHTTPS   string `json:"enable_https"`
}

// LoadConfiguration loads config from flags or .env file
func LoadConfiguration() *Config {
	cfg := Config{
		TokenExpire: TokenExpire,
		SecretKey:   SecretKey,
	}

	regStringVar(&cfg.ServerAddress, "a", ServerAddress, "Server address")
	regStringVar(&cfg.BaseURL, "b", BaseURL, "Server base URL")
	regStringVar(&cfg.LogLevel, "c", LogLevel, "Server log level")
	regStringVar(&cfg.FilePath, "f", FileName, "Server file storage")
	regStringVar(&cfg.DBPath, "d", DBPath, "Server db path")
	regStringVar(&cfg.EnableHTTPS, "s", EnableHTTPS, "Enable https")

	flag.Parse()

	flagServerAddress := getStringFlag("a")
	flagBaseURL := getStringFlag("b")
	flagFilePath := getStringFlag("f")
	flagDataBaseURI := getStringFlag("d")
	flagEnableHTTPS := getStringFlag("s")

	err := env.Parse(&cfg)

	if err != nil {
		panic(err)
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
	if flagEnableHTTPS != EnableHTTPS {
		cfg.EnableHTTPS = flagEnableHTTPS
	}

	configJSON, err := getJSONConfig()
	if err != nil {
		return &cfg
	}
	if configJSON.ServerAddress != ServerAddress {
		cfg.ServerAddress = flagServerAddress
	}
	if configJSON.BaseURL != BaseURL {
		cfg.BaseURL = flagBaseURL
	}
	if configJSON.FilePath != FileName {
		cfg.FilePath = flagFilePath
	}
	if configJSON.DBPath != DBPath {
		cfg.DBPath = flagDataBaseURI
	}
	if configJSON.EnableHTTPS != EnableHTTPS {
		cfg.EnableHTTPS = flagEnableHTTPS
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

func getJSONConfig() (JSONConfig, error) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	fmt.Println(exPath)
	var config JSONConfig
	jsonFile, errOpen := os.Open("/shortener/config/env.json")
	if errOpen != nil {
		fmt.Println(errOpen)
		return config, errOpen
	}
	defer jsonFile.Close()
	byteValue, errRead := io.ReadAll(jsonFile)
	if errRead != nil {
		return config, errRead
	}
	errUnmarshal := json.Unmarshal(byteValue, &config)
	if errUnmarshal != nil {
		return config, errRead
	}
	return config, nil
}
