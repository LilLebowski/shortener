package config

import (
	"flag"
	"os"
	"strings"
)

type ShortenerConfiguration struct {
	ServerAddress string
	BaseURL       string
}

func LoadConfiguration() *ShortenerConfiguration {
	cfg := &ShortenerConfiguration{}
	regStringVar(&cfg.ServerAddress, "a", "localhost:8080", "Server address")
	regStringVar(&cfg.BaseURL, "b", "http://localhost:8080", "Server base URL")
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
