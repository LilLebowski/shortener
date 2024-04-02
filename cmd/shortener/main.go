package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/LilLebowski/shortener/config"
	"github.com/LilLebowski/shortener/internal/handlers"
	"github.com/LilLebowski/shortener/internal/utils"
)

func main() {
	cfg := config.LoadConfiguration()
	err := utils.Initialize(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	storageInstance := utils.NewStorage()
	err = utils.FillFromStorage(storageInstance, cfg.FilePath)
	if err != nil {
		panic(err)
	}
	router := handlers.SetupRouter(cfg.BaseURL, storageInstance)
	router.Use(
		gin.Recovery(),
		utils.LoggerMiddleware(utils.Log),
		utils.CustomCompression(),
	)
	utils.Log.Info("Running server", zap.String("address", cfg.ServerAddress))
	routerErr := router.Run(cfg.ServerAddress)
	if routerErr != nil {
		panic(routerErr)
	}
	err = utils.WriteFile(storageInstance, cfg.FilePath, cfg.BaseURL)
	if err != nil {
		panic(err)
	}
}
