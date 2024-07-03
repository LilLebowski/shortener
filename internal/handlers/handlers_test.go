package handlers

import (
	"encoding/json"
	"fmt"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"

	"github.com/LilLebowski/shortener/cmd/shortener/config"
	"github.com/LilLebowski/shortener/internal/middleware"
	"github.com/LilLebowski/shortener/internal/mock_storage"
	"github.com/LilLebowski/shortener/internal/models"
	"github.com/LilLebowski/shortener/internal/services/shortener"
	"github.com/LilLebowski/shortener/internal/utils"
)

func buildJWTString(config *config.Config, userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, middleware.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.TokenExpire)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func setupRouter(strg *shortener.Service, config *config.Config) *gin.Engine {
	handlerWithService := Init(strg, config)

	router := gin.Default()
	router.Use(
		gin.Recovery(),
		middleware.Logger(middleware.Log),
		middleware.Compression(),
		middleware.Authorization(config),
	)
	router.GET("/ping", handlerWithService.GetPingHandler)
	router.GET("/:urlID", handlerWithService.GetShortURLHandler)
	router.POST("/", handlerWithService.CreateShortURLHandler)
	router.POST("/api/shorten", handlerWithService.CreateShortURLHandlerJSON)
	router.POST("/api/shorten/batch", handlerWithService.CreateBatch)
	router.GET("/api/user/urls", handlerWithService.GetListByUserIDHandler)
	router.DELETE("/api/user/urls", handlerWithService.DeleteUserUrlsHandler)

	router.HandleMethodNotAllowed = true

	return router
}

func TestCreateShortURLHandler(t *testing.T) {
	cfg := config.LoadConfiguration()

	type want struct {
		code int
	}
	tests := []struct {
		name            string
		URL             string
		userID          string
		short           string
		mockResponseErr error
		want            want
	}{
		{
			name:            "GET 1. body doesn't consist of data",
			URL:             "",
			short:           "",
			userID:          "bdf8817b-3225-4a46-9358-aa091b3cb478",
			mockResponseErr: fmt.Errorf(""),
			want: want{
				code: 400,
			},
		},
		{
			name:            "GET 2. body consist of data",
			URL:             "https://ya.ru",
			userID:          "bdf8817b-3225-4a46-9358-aa091b3cb478",
			short:           shortener.GetShortURL("https://ya.ru"),
			mockResponseErr: nil,
			want: want{
				code: 201,
			},
		},
		{
			name:            "GET 3. body conflict",
			URL:             "https://ya.ru",
			userID:          "bdf8817b-3225-4a46-9358-aa091b3cb478",
			short:           shortener.GetShortURL("https://ya.ru"),
			mockResponseErr: utils.NewUniqueConstraintError(fmt.Errorf("conflict")),
			want: want{
				code: 409,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tc *testing.T) {
			middleware.Initialize("debug")

			token, _ := buildJWTString(cfg, test.userID)
			param := strings.NewReader(test.URL)
			rq := httptest.NewRequest(http.MethodPost, "/", param)
			rw := httptest.NewRecorder()
			newCookie := http.Cookie{Name: "userID", Value: token}
			rq.Header.Add("Cookie", newCookie.String())

			ctrl := gomock.NewController(tc)
			strg := mock_storage.NewMockRepository(ctrl)
			s := &shortener.Service{
				BaseURL: config.BaseURL,
				Storage: strg,
			}

			strg.
				EXPECT().
				Set(test.URL, test.short, test.userID).
				Return(test.mockResponseErr).
				AnyTimes()

			router := setupRouter(s, cfg)
			router.ServeHTTP(rw, rq)

			res := rw.Result()
			defer res.Body.Close()
			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}

func TestCreateShortURLHandlerJSON(t *testing.T) {
	cfg := config.LoadConfiguration()

	type want struct {
		code int
	}
	tests := []struct {
		name            string
		userID          string
		short           string
		body            models.URLCreate
		want            want
		mockResponseErr error
	}{
		{
			name:   "GET 1. URL doesn't consist of data",
			short:  "",
			userID: "bdf8817b-3225-4a46-9358-aa091b3cb478",
			body: models.URLCreate{
				URL: "",
			},
			want: want{
				code: 400,
			},
			mockResponseErr: fmt.Errorf(""),
		},
		{
			name:   "GET 2. body consist of data",
			userID: "bdf8817b-3225-4a46-9358-aa091b3cb478",
			short:  shortener.GetShortURL("https://ya.ru"),
			body: models.URLCreate{
				URL: "https://ya.ru",
			},
			want: want{
				code: 201,
			},
			mockResponseErr: nil,
		},
		{
			name:   "GET 3. body conflict",
			userID: "bdf8817b-3225-4a46-9358-aa091b3cb478",
			short:  shortener.GetShortURL("https://ya.ru"),
			body: models.URLCreate{
				URL: "https://ya.ru",
			},
			want: want{
				code: 409,
			},
			mockResponseErr: utils.NewUniqueConstraintError(fmt.Errorf("conflict")),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tc *testing.T) {
			middleware.Initialize("debug")

			token, _ := buildJWTString(cfg, test.userID)
			jsonBytes, _ := json.Marshal(test.body)
			URL := strings.NewReader(string(jsonBytes))
			rq := httptest.NewRequest(http.MethodPost, "/api/shorten", URL)
			rw := httptest.NewRecorder()
			newCookie := http.Cookie{Name: "userID", Value: token}
			rq.Header.Set("Content-Type", "application/json")
			rq.Header.Add("Cookie", newCookie.String())

			ctrl := gomock.NewController(tc)
			strg := mock_storage.NewMockRepository(ctrl)
			s := &shortener.Service{
				BaseURL: config.BaseURL,
				Storage: strg,
			}

			strg.
				EXPECT().
				Set(test.body.URL, test.short, test.userID).
				Return(test.mockResponseErr).
				AnyTimes()

			router := setupRouter(s, cfg)
			router.ServeHTTP(rw, rq)

			res := rw.Result()
			defer res.Body.Close()
			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}

func TestGetShortURLHandler(t *testing.T) {
	cfg := config.LoadConfiguration()

	type want struct {
		code int
	}
	tests := []struct {
		name            string
		userID          string
		urlID           string
		url             string
		want            want
		mockResponseErr error
	}{
		{
			name:   "POST 1. URL doesn't exist",
			userID: "bdf8817b-3225-4a46-9358-aa091b3cb478",
			urlID:  "notfound",
			url:    "",
			want: want{
				code: 400,
			},
			mockResponseErr: utils.NewNotFoundError("url not found", fmt.Errorf("")),
		},
		{
			name:   "POST 2. URL exist",
			userID: "bdf8817b-3225-4a46-9358-aa091b3cb478",
			urlID:  "found",
			url:    "https://ya.ru",
			want: want{
				code: 307,
			},
			mockResponseErr: nil,
		},
		{
			name:   "POST 3. URL deleted",
			userID: "bdf8817b-3225-4a46-9358-aa091b3cb478",
			urlID:  "deleted",
			url:    "deleted",
			want: want{
				code: 410,
			},
			mockResponseErr: utils.NewDeletedError("url is already deleted", fmt.Errorf("")),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tc *testing.T) {
			middleware.Initialize("debug")

			token, _ := buildJWTString(cfg, test.userID)

			ctrl := gomock.NewController(tc)
			strg := mock_storage.NewMockRepository(ctrl)
			s := &shortener.Service{
				BaseURL: config.BaseURL,
				Storage: strg,
			}

			rq := httptest.NewRequest(http.MethodGet, "/"+test.urlID, nil)
			rw := httptest.NewRecorder()
			newCookie := http.Cookie{Name: "userID", Value: token}
			rq.Header.Set("Content-Type", "application/json")
			rq.Header.Add("Cookie", newCookie.String())

			strg.
				EXPECT().
				Get(test.urlID).
				Return(test.url, test.mockResponseErr).
				AnyTimes()

			router := setupRouter(s, cfg)
			router.ServeHTTP(rw, rq)

			res := rw.Result()
			defer res.Body.Close()
			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}

func TestCreateBatch(t *testing.T) {
	cfg := config.LoadConfiguration()

	type args struct {
		code        int
		contentType string
	}
	tests := []struct {
		name            string
		userID          string
		args            args
		body            []models.URLs
		mock            []models.FullURLs
		mockResponseErr error
	}{
		{
			name:   "POST 1. Full body",
			userID: "bdf8817b-3225-4a46-9358-aa091b3cb478",
			args: args{
				code:        201,
				contentType: "application/json",
			},
			body: []models.URLs{
				{
					CorrelationID: "1",
					OriginalURL:   "https://ya.ru",
				},
				{
					CorrelationID: "2",
					OriginalURL:   "https://ya.com",
				},
			},
			mock: []models.FullURLs{
				{
					OriginalURL: "https://ya.ru",
					ShortURL:    shortener.GetShortURL("https://ya.ru"),
				},
				{
					OriginalURL: "https://ya.com",
					ShortURL:    shortener.GetShortURL("https://ya.com"),
				},
			},
			mockResponseErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tc *testing.T) {
			middleware.Initialize("debug")

			token, _ := buildJWTString(cfg, test.userID)

			ctrl := gomock.NewController(tc)
			strg := mock_storage.NewMockRepository(ctrl)
			s := &shortener.Service{
				BaseURL: config.BaseURL,
				Storage: strg,
			}

			jsonBody, err := json.Marshal(test.body)
			if err != nil {
				t.Fatal(err)
			}
			rq := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(string(jsonBody)))
			rw := httptest.NewRecorder()
			newCookie := http.Cookie{Name: "userID", Value: token}
			rq.Header.Set("Content-Type", "application/json")
			rq.Header.Add("Cookie", newCookie.String())

			strg.
				EXPECT().
				SetBatch(test.userID, test.mock).
				Return(test.mockResponseErr).
				AnyTimes()

			router := setupRouter(s, cfg)
			router.ServeHTTP(rw, rq)

			res := rw.Result()
			defer res.Body.Close()

			assert.Equal(t, test.args.code, res.StatusCode)
			assert.Equal(t, test.args.contentType, res.Header.Get("Content-Type"))
		})
	}
}
