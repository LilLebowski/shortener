package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"go.uber.org/zap"

	"github.com/LilLebowski/shortener/cmd/shortener/config"
	"github.com/LilLebowski/shortener/internal/middleware"
	"github.com/LilLebowski/shortener/internal/router"
)

// Global variables
var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	logBuildInfo()

	ctx, cancel := context.WithCancel(context.Background())

	cfg := config.LoadConfiguration()
	err := middleware.Initialize(cfg.LogLevel)
	if err != nil {
		panic(err)
	}

	routerInstance := router.Init(cfg)

	middleware.Log.Info("Running server", zap.String("address", cfg.ServerAddress))
	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: routerInstance,
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

// logBuildInfo print info about package
func logBuildInfo() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}

	if buildDate == "" {
		buildDate = "N/A"
	}

	if buildCommit == "" {
		buildCommit = "N/A"
	}

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
