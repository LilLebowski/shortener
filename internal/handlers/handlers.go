package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

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

func SetupRouter(configBaseURL string, storageInstance *storage.Storage) *gin.Engine {
	storageShortener := shortener.Init(configBaseURL, storageInstance)

	router := gin.Default()
	router.Use(
		gin.Recovery(),
		utils.LoggerMiddleware(utils.Log),
		utils.CustomCompression(),
	)
	router.GET("/ping", GetPingHandler(storageShortener))
	router.GET("/:urlID", GetShortURLHandler(storageShortener))
	router.POST("/", CreateShortURLHandler(storageShortener))
	router.POST("/api/shorten", CreateShortURLHandlerJSON(storageShortener))
	router.POST("/api/shorten/batch", CreateBatch(storageShortener))

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
		shortURL, setErr := sh.Set(url)
		if setErr != nil {
			var uce *utils.UniqueConstraintError
			if errors.As(setErr, &uce) {
				ctx.Writer.WriteHeader(http.StatusConflict)
			}
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
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
		if value, ok := sh.Get(urlID); ok {
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
		shortURL, setErr := sh.Set(reqBody.URL)
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
			utils.Log.Fatalf("cannot write response to the client: %s", err)
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
			shortURL, setErr := sh.Set(url)
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
