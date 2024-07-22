package main

import (
	"context"
	"fmt"
	"github.com/LilLebowski/shortener/internal/services/shortener"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LilLebowski/shortener/cmd/shortener/config"
	"github.com/LilLebowski/shortener/internal/middleware"
	"github.com/LilLebowski/shortener/internal/router"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
)

// Global variables
const (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	logBuildInfo()

	cfg := config.LoadConfiguration()
	err := middleware.Initialize(cfg.LogLevel)
	if err != nil {
		panic(err)
	}

	routerInstance, service := router.Init(cfg)

	middleware.Log.Info("Running server", zap.String("address", cfg.ServerAddress))

	ctx := context.Background()

	if cfg.EnableHTTPS == "" {
		srv := startHTTPServer(cfg, routerInstance)
		shutDownServer(ctx, cfg, srv, service)
	} else {
		srv := startHTTPSServer(cfg, routerInstance)
		shutDownServer(ctx, cfg, srv, service)
	}
}

// startHTTPSServer run HTTPS server
func startHTTPSServer(c *config.Config, r *gin.Engine) *http.Server {
	manager := &autocert.Manager{
		Cache:      autocert.DirCache("cache-dir"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(c.ServerAddress),
	}

	srv := &http.Server{
		Addr:      c.ServerAddress,
		Handler:   r,
		TLSConfig: manager.TLSConfig(),
	}

	go func() {
		err := srv.ListenAndServeTLS("server.crt", "server.key")
		if err != nil {
			middleware.Log.Info("app error exit", zap.Error(err))
			syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		}
	}()

	middleware.Log.Info("Running server", zap.String("address", c.ServerAddress))

	return srv
}

// startHTTPServer run HTTP server
func startHTTPServer(c *config.Config, r *gin.Engine) *http.Server {
	srv := &http.Server{
		Addr:    c.ServerAddress,
		Handler: r,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			middleware.Log.Info("app error exit", zap.Error(err))
			syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		}
	}()

	middleware.Log.Info("Running server", zap.String("address", c.ServerAddress))

	return srv
}

// shutDownServer exit for server
func shutDownServer(ctx context.Context, c *config.Config, srv *http.Server, sh *shortener.Service) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	killSignal := <-interrupt
	switch killSignal {
	case os.Interrupt:
		middleware.Log.Info("Got SIGINT...")
	case syscall.SIGTERM:
		middleware.Log.Info("Got SIGTERM...")
	}

	releaseResources(ctx, c, srv, sh)
}

// releaseResources free resources
func releaseResources(ctx context.Context, c *config.Config, srv *http.Server, sh *shortener.Service) {
	middleware.Log.Info("The service is shutting down...")
	if c.DBPath != "" {
		middleware.Log.Info("Closing connect to db")
		err := sh.Storage.Close()
		if err != nil {
			middleware.Log.Info("Error while closing db connection", zap.Error(err))
		}
	}
	time.Sleep(1 * time.Second)
	if err := srv.Shutdown(ctx); err != nil {
		middleware.Log.Info("app exit error", zap.Error(err))
	}
	middleware.Log.Info("Done")
}

// logBuildInfo print info about package
func logBuildInfo() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
