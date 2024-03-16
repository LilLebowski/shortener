package main

import (
	"fmt"

	"github.com/LilLebowski/shortener/config"
	"github.com/LilLebowski/shortener/internal/handlers"
)

func main() {
	cfg := config.LoadConfiguration()
	router := handlers.SetupRouter(cfg.BaseURL)
	fmt.Printf("Server Address: %s\n", cfg.ServerAddress)
	routerErr := router.Run(cfg.ServerAddress)
	if routerErr != nil {
		panic(routerErr)
	}
}
