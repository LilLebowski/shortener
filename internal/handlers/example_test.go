package handlers

import (
	"encoding/json"
	"github.com/LilLebowski/shortener/cmd/shortener/config"
	"github.com/LilLebowski/shortener/internal/middleware"
	"github.com/LilLebowski/shortener/internal/models"
	"github.com/LilLebowski/shortener/internal/services/shortener"
	fs "github.com/LilLebowski/shortener/internal/storage/file"
	"net/http"
	"net/http/httptest"
	"strings"
)

func Example_createShortURLHandler() {
	cfg := config.LoadConfiguration()
	middleware.Initialize("debug")

	userID := "bdf8817b-3225-4a46-9358-aa091b3cb478"
	URL := "https://ya.ru"
	token, _ := buildJWTString(cfg, userID)

	param := strings.NewReader(URL)
	rq := httptest.NewRequest(http.MethodPost, "/", param)
	rw := httptest.NewRecorder()
	newCookie := http.Cookie{Name: "userID", Value: token}
	rq.Header.Add("Cookie", newCookie.String())

	strg := fs.Init("")
	s := &shortener.Service{
		BaseURL: config.BaseURL,
		Storage: strg,
	}
	router := setupRouter(s, cfg)

	router.ServeHTTP(rw, rq)
	res := rw.Result()

	defer res.Body.Close()
}

func Example_createShortURLHandlerJSON() {
	cfg := config.LoadConfiguration()
	middleware.Initialize("debug")

	body := models.URLCreate{
		URL: "https://ya.ru",
	}
	userID := "bdf8817b-3225-4a46-9358-aa091b3cb478"

	strg := fs.Init("")
	s := &shortener.Service{
		BaseURL: config.BaseURL,
		Storage: strg,
	}

	token, _ := buildJWTString(cfg, userID)
	jsonBytes, _ := json.Marshal(body)
	URL := strings.NewReader(string(jsonBytes))
	rq := httptest.NewRequest(http.MethodPost, "/api/shorten", URL)
	rw := httptest.NewRecorder()
	newCookie := http.Cookie{Name: "userID", Value: token}
	rq.Header.Set("Content-Type", "application/json")
	rq.Header.Add("Cookie", newCookie.String())

	router := setupRouter(s, cfg)
	router.ServeHTTP(rw, rq)
	res := rw.Result()

	defer res.Body.Close()
}

func Example_getShortURLHandler() {
	cfg := config.LoadConfiguration()
	middleware.Initialize("debug")

	userID := "bdf8817b-3225-4a46-9358-aa091b3cb478"
	urlID := "found"
	token, _ := buildJWTString(cfg, userID)

	strg := fs.Init("")
	s := &shortener.Service{
		BaseURL: config.BaseURL,
		Storage: strg,
	}

	rq := httptest.NewRequest(http.MethodGet, "/"+urlID, nil)
	rw := httptest.NewRecorder()
	newCookie := http.Cookie{Name: "userID", Value: token}
	rq.Header.Set("Content-Type", "application/json")
	rq.Header.Add("Cookie", newCookie.String())

	router := setupRouter(s, cfg)

	router.ServeHTTP(rw, rq)
	res := rw.Result()

	defer res.Body.Close()
}

func Example_shortenURLsHandlerJSON() {
	cfg := config.LoadConfiguration()
	middleware.Initialize("debug")

	body := []models.URLs{
		{
			CorrelationID: "1",
			OriginalURL:   "https://ya.ru",
		},
		{
			CorrelationID: "2",
			OriginalURL:   "https://ya.com",
		},
	}
	userID := "bdf8817b-3225-4a46-9358-aa091b3cb478"
	strg := fs.Init("")
	s := &shortener.Service{
		BaseURL: config.BaseURL,
		Storage: strg,
	}
	token, _ := buildJWTString(cfg, userID)

	jsonBody, _ := json.Marshal(body)
	rq := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(string(jsonBody)))
	rw := httptest.NewRecorder()
	newCookie := http.Cookie{Name: "userID", Value: token}
	rq.Header.Set("Content-Type", "application/json")
	rq.Header.Add("Cookie", newCookie.String())

	router := setupRouter(s, cfg)

	router.ServeHTTP(rw, rq)
	res := rw.Result()

	defer res.Body.Close()
}
