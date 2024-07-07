// Package router contains general routes for shortener service
package router

import (
	"github.com/gin-gonic/gin"

	"github.com/LilLebowski/shortener/cmd/shortener/config"
	"github.com/LilLebowski/shortener/internal/handlers"
	"github.com/LilLebowski/shortener/internal/middleware"
	"github.com/LilLebowski/shortener/internal/services/shortener"
)

// Init define routes
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
	// create link from text format
	router.POST("/", handlerWithService.CreateShortURLHandler)
	// create link from JSON format
	router.POST("/api/shorten", handlerWithService.CreateShortURLHandlerJSON)
	// get original url by urlID
	router.GET("/:urlID", handlerWithService.GetShortURLHandler)
	// batch creation short links
	router.POST("/api/shorten/batch", handlerWithService.CreateBatch)
	// get user session links in JSON
	router.GET("/api/user/urls", handlerWithService.GetListByUserIDHandler)
	// delete links session
	router.DELETE("/api/user/urls", handlerWithService.DeleteUserUrlsHandler)
	// ping db
	router.GET("/ping", handlerWithService.GetPingHandler)

	router.HandleMethodNotAllowed = true

	return router
}
