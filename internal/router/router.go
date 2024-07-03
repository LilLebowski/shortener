package router

import (
	"github.com/gin-gonic/gin"

	"github.com/LilLebowski/shortener/cmd/shortener/config"
	"github.com/LilLebowski/shortener/internal/handlers"
	"github.com/LilLebowski/shortener/internal/middleware"
	"github.com/LilLebowski/shortener/internal/services/shortener"
)

func Init(config *config.Config) *gin.Engine {
	shortenerService := shortener.Init(config)
	handlerWithService := handlers.Init(shortenerService, config)

	router := gin.Default()
	router.Use(
		gin.Recovery(),
		middleware.Logger(middleware.Log),
		middleware.Compression(),
		middleware.Authorization(config),
	)
	router.POST("/", handlerWithService.CreateShortURLHandler)
	router.POST("/api/shorten", handlerWithService.CreateShortURLHandlerJSON)
	router.GET("/:urlID", handlerWithService.GetShortURLHandler)
	router.POST("/api/shorten/batch", handlerWithService.CreateBatch)
	router.GET("/api/user/urls", handlerWithService.GetListByUserIDHandler)
	router.DELETE("/api/user/urls", handlerWithService.DeleteUserUrlsHandler)
	router.GET("/ping", handlerWithService.GetPingHandler)

	router.HandleMethodNotAllowed = true

	return router
}
