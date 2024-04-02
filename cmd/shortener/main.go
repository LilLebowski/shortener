package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/gzip"
	"go.uber.org/zap"

	"github.com/LilLebowski/shortener/config"
	"github.com/LilLebowski/shortener/internal/handlers"
	"github.com/LilLebowski/shortener/internal/utils"
)

func main() {
	cfg := config.LoadConfiguration()
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("logger don't Run! %s", err)
	}
	utils.Sugar = zapLogger.Sugar()
	router := handlers.SetupRouter(cfg.BaseURL)
	router.Use(utils.CustomMiddlewareLogger)
	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(utils.CustomCompression)))
	fmt.Printf("Server Address: %s\n", cfg.ServerAddress)
	routerErr := router.Run(cfg.ServerAddress)
	if routerErr != nil {
		panic(routerErr)
	}
}
