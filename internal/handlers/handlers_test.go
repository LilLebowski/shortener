package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/LilLebowski/shortener/cmd/shortener/config"
	"github.com/LilLebowski/shortener/internal/storage"
	"github.com/LilLebowski/shortener/internal/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateShortURLHandler(t *testing.T) {
	cfg := config.LoadConfiguration()

	type want struct {
		code int
	}
	tests := []struct {
		name  string
		param string
		want  want
	}{
		{
			name:  "GET 1. body doesn't consist of data",
			param: "",
			want: want{
				code: 400,
			},
		},
		{
			name:  "GET 2. body consist of data",
			param: "https://ya.ru",
			want: want{
				code: 201,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Printf("\n\nTest %v Body %v\n", cfg.BaseURL, test.param)
			storageInstance := storage.Init(cfg.FilePath, cfg.DBPath)
			utils.Initialize("debug")
			router := SetupRouter(cfg.BaseURL, storageInstance)
			param := strings.NewReader(test.param)
			rq := httptest.NewRequest(http.MethodPost, "/", param)
			rw := httptest.NewRecorder()
			router.ServeHTTP(rw, rq)
			res := rw.Result()
			defer res.Body.Close()
			fmt.Printf("want code = %d StatusCode %d\n", test.want.code, res.StatusCode)
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
		name  string
		urlID string
		url   string
		want  want
	}{
		{
			name:  "POST 1. URL doesn't exist",
			urlID: "notfound",
			url:   "notfound",
			want: want{
				code: 400,
			},
		},
		{
			name:  "POST 2. URL exist",
			urlID: "found",
			url:   "https://ya.ru",
			want: want{
				code: 307,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Printf("\n\nTest %v urlID %v url %v\n", test.name, test.urlID, test.url)
			storageInstance := storage.Init(cfg.FilePath, cfg.DBPath)
			utils.Initialize("debug")
			router := SetupRouter(cfg.BaseURL, storageInstance)
			if test.urlID == "found" {
				storageInstance.Memory.Set(test.url, test.urlID)
			}
			rq := httptest.NewRequest(http.MethodGet, "/"+test.urlID, nil)
			rw := httptest.NewRecorder()
			router.ServeHTTP(rw, rq)
			res := rw.Result()
			defer res.Body.Close()
			fmt.Printf("want code = %d StatusCode %d\n", test.want.code, res.StatusCode)
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
		name string
		body CreateURLData
		want want
	}{
		{
			name: "GET 1. URL doesn't consist of data",
			body: CreateURLData{
				URL: "",
			},
			want: want{
				code: 400,
			},
		},
		{
			name: "GET 2. body consist of data",
			body: CreateURLData{
				URL: "https://ya.ru",
			},
			want: want{
				code: 201,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Printf("\n\nTest %v Body %v\n", cfg.BaseURL, test.body)
			storageInstance := storage.Init(cfg.FilePath, cfg.DBPath)
			utils.Initialize("debug")
			router := SetupRouter(cfg.BaseURL, storageInstance)
			jsonBytes, _ := json.Marshal(test.body)
			param := strings.NewReader(string(jsonBytes))
			rq := httptest.NewRequest(http.MethodPost, "/api/shorten", param)
			rq.Header.Set("Content-Type", "application/json")
			rw := httptest.NewRecorder()
			router.ServeHTTP(rw, rq)
			res := rw.Result()
			defer res.Body.Close()
			fmt.Printf("want code = %d StatusCode %d\n", test.want.code, res.StatusCode)
			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}
