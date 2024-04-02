package main

import (
	"fmt"
	"log"

	"go.uber.org/zap"

	"github.com/LilLebowski/shortener/config"
	"github.com/LilLebowski/shortener/internal/handlers"
	"github.com/LilLebowski/shortener/internal/logger"
)

func main() {
	cfg := config.LoadConfiguration()
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("logger don't Run! %s", err)
	}
	logger.Sugar = zapLogger.Sugar()
	router := handlers.SetupRouter(cfg.BaseURL)
	router.Use(logger.CustomMiddlewareLogger())
	fmt.Printf("Server Address: %s\n", cfg.ServerAddress)
	routerErr := router.Run(cfg.ServerAddress)
	if routerErr != nil {
		panic(routerErr)
	}
}
