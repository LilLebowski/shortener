package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/LilLebowski/shortener/internal/utils"
	"github.com/gin-gonic/gin"
)

var urls map[string]string
var baseURL string

func SetupRouter(configBaseURL string) *gin.Engine {
	urls = make(map[string]string)
	baseURL = configBaseURL
	fmt.Printf("base URL: %s\n", baseURL)

	router := gin.Default()
	router.GET("/:urlID", GetShortURLHandler)
	router.POST("/", CreateShortURLHandler)
	router.POST("/api/shorten", CreateShortURLHandlerJSON)

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
	fmt.Printf("request body: %s\n", reqBodyString)
	if reqBodyString == "" {
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	res, encodeErr := utils.EncodeURL(reqBodyString)
	if encodeErr != nil {
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	urls[res] = reqBodyString
	ctx.Writer.Header().Set("Content-Type", "text/plain")
	ctx.Writer.WriteHeader(http.StatusCreated)
	newAddr := baseURL + "/" + res
	_, writeErr := ctx.Writer.Write([]byte(newAddr))
	if writeErr != nil {
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

type CreateURLData struct {
	URL string `json:"URL"`
}

type CreateURLResponse struct {
	Result string `json:"result"`
}

func CreateShortURLHandlerJSON(ctx *gin.Context) {
	if ctx.Request.Header.Get("Content-Type") != "application/json" {
		http.Error(ctx.Writer, "Invalid Content Type!", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		http.Error(ctx.Writer, fmt.Sprintf("Cannot read request body: %s", err), http.StatusBadRequest)
		return
	}

	var reqBody CreateURLData
	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		http.Error(ctx.Writer, fmt.Sprintf("Cannot decode request body to `JSON`: %s", err), http.StatusBadRequest)
		return
	}
	res, encodeErr := utils.EncodeURL(reqBody.URL)
	if encodeErr != nil {
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if reqBody.URL == "" {
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("request body: %s\n", reqBody)
	urls[res] = reqBody.URL
	ctx.Writer.Header().Set("Content-Type", "application/json")
	ctx.Writer.WriteHeader(http.StatusCreated)
	shortRes := CreateURLResponse{
		Result: baseURL + "/" + res,
	}
	resp, err := json.Marshal(shortRes)
	if err != nil {
		http.Error(ctx.Writer, fmt.Sprintf("cannot encode response: %s", err), http.StatusBadRequest)
	}

	_, err = ctx.Writer.Write(resp)
	if err != nil {
		log.Fatalf("cannot write response to the client: %s", err)
	}
}
