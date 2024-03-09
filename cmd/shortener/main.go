package main

import (
	"github.com/LilLebowski/shortener/config"
	"github.com/LilLebowski/shortener/internal/handlers"
)

func main() {
	cfg := config.LoadConfiguration()
	router := handlers.SetupRouter(cfg.BaseURL)
	routerErr := router.Run(cfg.ServerAddress)
	if routerErr != nil {
		panic(routerErr)
	}
}
