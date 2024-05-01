package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/LilLebowski/shortener/cmd/shortener/config"
	"github.com/LilLebowski/shortener/internal/middleware"
	"github.com/LilLebowski/shortener/internal/services/shortener"
	"github.com/LilLebowski/shortener/internal/storage"
	"github.com/LilLebowski/shortener/internal/utils"
)

type CreateBatchData struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type CreateBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type CreateURLData struct {
	URL string `json:"URL"`
}

type CreateURLResponse struct {
	Result string `json:"result"`
}

func SetupRouter(config *config.Config, storageInstance *storage.Storage) *gin.Engine {
	storageShortener := shortener.Init(config.BaseURL, storageInstance)

	router := gin.Default()
	router.Use(
		gin.Recovery(),
		middleware.Logger(middleware.Log),
		middleware.Compression(),
		middleware.Authorization(config),
	)
	router.GET("/ping", GetPingHandler(storageShortener))
	router.GET("/:urlID", GetShortURLHandler(storageShortener))
	router.POST("/", CreateShortURLHandler(storageShortener))
	router.POST("/api/shorten", CreateShortURLHandlerJSON(storageShortener))
	router.POST("/api/shorten/batch", CreateBatch(storageShortener))
	router.GET("/api/user/urls", GetListByUserIDHandler(storageShortener))
	router.DELETE("/api/user/urls", DeleteUserUrlsHandler(storageShortener))

	router.HandleMethodNotAllowed = true

	return router
}

func CreateShortURLHandler(sh *shortener.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqBody, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			fmt.Printf("could not read request body: %s\n", err)
			ctx.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		url := strings.TrimSpace(string(reqBody))
		fmt.Printf("request url: %s\n", url)
		if url == "" {
			ctx.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		userIDFromContext, _ := ctx.Get("userID")
		userID, _ := userIDFromContext.(string)
		shortURL, setErr := sh.Set(url, userID)
		if setErr != nil {
			var uce *utils.UniqueConstraintError
			if errors.As(setErr, &uce) {
				ctx.Writer.WriteHeader(http.StatusConflict)
			} else {
				ctx.Writer.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			ctx.Writer.WriteHeader(http.StatusCreated)
		}
		ctx.Writer.Header().Set("Content-Type", "text/plain")
		_, writeErr := ctx.Writer.Write([]byte(shortURL))
		if writeErr != nil {
			ctx.Writer.WriteHeader(http.StatusBadRequest)
		}
	}
}

func GetShortURLHandler(sh *shortener.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		urlID := ctx.Param("urlID")
		fmt.Printf("url id: %s\n", urlID)
		value, isDeleted, ok := sh.Get(urlID)
		if ok {
			if isDeleted {
				ctx.Status(http.StatusGone)
				return
			}
			fmt.Printf("found url: %s\n", value)
			ctx.Writer.Header().Set("Location", value)
			ctx.Writer.WriteHeader(http.StatusTemporaryRedirect)
		} else {
			ctx.Writer.Header().Set("Location", "Not found")
			ctx.Writer.WriteHeader(http.StatusBadRequest)
		}
	}
}

func CreateShortURLHandlerJSON(sh *shortener.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
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
		if reqBody.URL == "" {
			ctx.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Printf("request body: %s\n", reqBody)
		userIDFromContext, _ := ctx.Get("userID")
		userID, _ := userIDFromContext.(string)
		shortURL, setErr := sh.Set(reqBody.URL, userID)
		if setErr != nil {
			var uce *utils.UniqueConstraintError
			if errors.As(setErr, &uce) {
				ctx.Writer.WriteHeader(http.StatusConflict)
			} else {
				ctx.Writer.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			ctx.Writer.WriteHeader(http.StatusCreated)
		}
		ctx.Writer.Header().Set("Content-Type", "application/json")
		shortRes := CreateURLResponse{
			Result: shortURL,
		}
		resp, err := json.Marshal(shortRes)
		if err != nil {
			http.Error(ctx.Writer, fmt.Sprintf("cannot encode response: %s", err), http.StatusBadRequest)
		}

		_, err = ctx.Writer.Write(resp)
		if err != nil {
			middleware.Log.Fatalf("cannot write response to the client: %s", err)
		}
	}
}

func GetPingHandler(sh *shortener.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := sh.Ping()
		if err != nil {
			fmt.Printf("err: %s", err)
			ctx.JSON(http.StatusInternalServerError, "")
			return
		}
		ctx.JSON(http.StatusOK, "")
	}
}

func CreateBatch(sh *shortener.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var decoderBody []CreateBatchData
		decoder := json.NewDecoder(ctx.Request.Body)
		err := decoder.Decode(&decoderBody)
		if err != nil {
			errorMassage := map[string]interface{}{
				"message": "Failed to read request body",
				"code":    http.StatusInternalServerError,
			}
			answer, _ := json.Marshal(errorMassage)
			ctx.Data(http.StatusInternalServerError, "application/json", answer)
			return
		}

		httpStatus := http.StatusCreated
		var URLResponses []CreateBatchResponse
		for _, req := range decoderBody {
			url := strings.TrimSpace(req.OriginalURL)
			userIDFromContext, _ := ctx.Get("userID")
			userID, _ := userIDFromContext.(string)
			shortURL, setErr := sh.Set(url, userID)
			if setErr != nil {
				var uce *utils.UniqueConstraintError
				if errors.As(setErr, &uce) {
					httpStatus = http.StatusConflict
				} else {
					errorMassage := map[string]interface{}{
						"message": "the url could not be shortened",
						"code":    http.StatusInternalServerError,
					}
					answer, _ := json.Marshal(errorMassage)
					ctx.Data(http.StatusInternalServerError, "application/json", answer)
					return
				}
			}
			urlResponse := CreateBatchResponse{
				req.CorrelationID,
				shortURL,
			}
			URLResponses = append(URLResponses, urlResponse)
		}

		respJSON, err := json.Marshal(URLResponses)
		if err != nil {
			errorMassage := map[string]interface{}{
				"message": "Failed to read request body",
				"code":    http.StatusInternalServerError,
			}
			answer, _ := json.Marshal(errorMassage)
			ctx.Data(http.StatusInternalServerError, "application/json", answer)
			return
		}
		ctx.Data(httpStatus, "application/json", respJSON)
	}
}

func GetListByUserIDHandler(sh *shortener.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		code := http.StatusOK
		userIDFromContext, exists := ctx.Get("userID")
		if !exists {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to get userID",
				"error":   errors.New("failed to get user from context").Error(),
			})
			return
		}
		isNew, _ := ctx.Get("new")
		if isNew == true {
			code = http.StatusUnauthorized
			ctx.JSON(code, nil)
			return
		}
		userID, _ := userIDFromContext.(string)
		urls, err := sh.GetByUserID(userID)
		ctx.Header("Content-type", "application/json")
		if err != nil {
			code = http.StatusInternalServerError
			ctx.JSON(code, gin.H{
				"message": "Failed to retrieve user URLs",
				"code":    code,
			})
			return
		}

		if len(urls) == 0 {
			ctx.JSON(http.StatusNoContent, nil)
			return
		}
		ctx.JSON(code, urls)
	}
}

func DeleteUserUrlsHandler(sh *shortener.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		code := http.StatusAccepted
		userIDFromContext, exists := ctx.Get("userID")
		if !exists {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to get userID",
				"error":   errors.New("failed to get user from context").Error(),
			})
			return
		}
		userID, _ := userIDFromContext.(string)

		var shorURLs []string
		if err := ctx.BindJSON(&shorURLs); err != nil {
			code = http.StatusBadRequest
			ctx.JSON(code, gin.H{
				"error:": err.Error(),
			})
		}

		err := sh.DeleteURLsRep(userID, shorURLs)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed delete to url",
				"error":   errors.New("failed to get user from context").Error(),
			})
			return
		}
		ctx.Status(code)
	}
}
