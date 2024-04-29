package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"go.uber.org/zap"

	"github.com/LilLebowski/shortener/cmd/shortener/config"
	"github.com/LilLebowski/shortener/internal/handlers"
	"github.com/LilLebowski/shortener/internal/storage"
	"github.com/LilLebowski/shortener/internal/utils"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cfg := config.LoadConfiguration()
	err := utils.Initialize(cfg.LogLevel)
	if err != nil {
		panic(err)
	}

	storageInstance := storage.Init(cfg.FilePath, cfg.DBPath)

	handler := handlers.SetupRouter(cfg.BaseURL, storageInstance)

	utils.Log.Info("Running server", zap.String("address", cfg.ServerAddress))
	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: handler,
	}

	go func() {
		log.Println(server.ListenAndServe())
		cancel()
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	select {
	case <-sigint:
		cancel()
	case <-ctx.Done():
	}
	err = server.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}
}
