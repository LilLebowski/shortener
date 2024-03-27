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
	regStringVar(&cfg.ServerAddress, "a", "localhost:8080", "Server address")
	regStringVar(&cfg.BaseURL, "b", "http://localhost:8080", "Server base URL")
	flag.Parse()
	if cfg.ServerAddress = os.Getenv("SERVER_ADDRESS"); cfg.ServerAddress == "" {
		cfg.ServerAddress = getStringFlag("a")
	}

	if cfg.BaseURL = os.Getenv("BASE_URL"); cfg.BaseURL == "" {
		cfg.BaseURL = getStringFlag("b")
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
