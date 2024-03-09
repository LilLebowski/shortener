package handlers

import (
	"fmt"
	"github.com/LilLebowski/shortener/internal/utils"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

var urls map[string]string
var baseURL string

func SetupRouter(configBaseURL string) *gin.Engine {
	urls = make(map[string]string)
	baseURL = configBaseURL

	router := gin.Default()
	router.GET("/:urlID", GetShortURLHandler)
	router.POST("/", CreateShortURLHandler)

	return router
}

func CreateShortURLHandler(ctx *gin.Context) {
	reqBody, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		fmt.Printf("could not read request body: %s\n", err)
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	reqBodyString := string(reqBody)
	fmt.Printf("request body: %s\n", baseURL)
	if reqBodyString != "" {
		res, encodeErr := utils.EncodeURL(reqBodyString)
		if encodeErr == nil {
			urls[res] = reqBodyString
			ctx.Writer.Header().Set("Content-Type", "text/plain")
			ctx.Writer.WriteHeader(http.StatusCreated)
			newAddr := baseURL + res
			_, writeErr := ctx.Writer.Write([]byte(newAddr))
			if writeErr != nil {
				ctx.Writer.WriteHeader(http.StatusBadRequest)
			}
		} else {
			ctx.Writer.WriteHeader(http.StatusBadRequest)
		}
	} else {
		ctx.Writer.WriteHeader(http.StatusBadRequest)
	}
}

func GetShortURLHandler(ctx *gin.Context) {
	fmt.Printf("current session: %s\n", urls)
	urlID := ctx.Param("urlID")
	fmt.Printf("url id: %s\n", urlID)
	if value, ok := urls[urlID]; ok {
		fmt.Printf("found url: %s\n", value)
		ctx.Writer.Header().Set("Location", value)
		ctx.Writer.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		ctx.Writer.Header().Set("Location", "Not found")
		ctx.Writer.WriteHeader(http.StatusBadRequest)
	}
}
