package utils

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var Sugar *zap.SugaredLogger

func CustomMiddlewareLogger(ctx *gin.Context) {
	start := time.Now()
	Sugar.Infoln("URI:", ctx.Request.RequestURI, "Method:", ctx.Request.Method, "Duration:", time.Since(start))
}
