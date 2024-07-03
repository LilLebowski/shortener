package handlers

//
//import (
//	"encoding/json"
//	"net/http"
//	"net/http/httptest"
//	"strconv"
//	"strings"
//	"testing"
//
//	"github.com/LilLebowski/shortener/cmd/shortener/config"
//	"github.com/LilLebowski/shortener/internal/middleware"
//	"github.com/LilLebowski/shortener/internal/storage"
//)
//
//func BenchmarkCreateShortURLHandler(b *testing.B) {
//	cfg := config.LoadConfiguration()
//	storageInstance := storage.Init(cfg.FilePath, cfg.DBPath)
//	middleware.Initialize("debug")
//	router := setupRouter(cfg, storageInstance)
//
//	b.ResetTimer()
//
//	for i := 0; i < b.N; i++ {
//		b.StopTimer() // stop all timers
//		URL := strings.NewReader("https://ya.ru/" + strconv.Itoa(i))
//		rq := httptest.NewRequest(http.MethodPost, "/", URL)
//		rw := httptest.NewRecorder()
//		b.StartTimer()
//		router.ServeHTTP(rw, rq)
//		res := rw.Result()
//		b.StopTimer()
//		defer res.Body.Close()
//	}
//}
//
//func BenchmarkGetShortURLHandler(b *testing.B) {
//	cfg := config.LoadConfiguration()
//	storageInstance := storage.Init(cfg.FilePath, cfg.DBPath)
//	middleware.Initialize("debug")
//	router := setupRouter(cfg, storageInstance)
//
//	b.ResetTimer()
//
//	for i := 0; i < b.N; i++ {
//		storageInstance.Memory.Set("https://ya.ru/"+strconv.Itoa(i), "urlID", "")
//		b.StopTimer() // stop all timers
//		rq := httptest.NewRequest(http.MethodGet, "/urlID", nil)
//		rw := httptest.NewRecorder()
//		b.StartTimer()
//		router.ServeHTTP(rw, rq)
//		res := rw.Result()
//		b.StopTimer()
//		defer res.Body.Close()
//	}
//}
//
//func BenchmarkCreateShortURLHandlerJSON(b *testing.B) {
//	cfg := config.LoadConfiguration()
//	storageInstance := storage.Init(cfg.FilePath, cfg.DBPath)
//	middleware.Initialize("debug")
//	router := setupRouter(cfg, storageInstance)
//
//	b.ResetTimer()
//
//	for i := 0; i < b.N; i++ {
//		b.StopTimer() // stop all timers
//		body := CreateURLData{
//			URL: "https://ya.ru/" + strconv.Itoa(i),
//		}
//		jsonBytes, _ := json.Marshal(body)
//		URL := strings.NewReader(string(jsonBytes))
//		rq := httptest.NewRequest(http.MethodPost, "/api/shorten", URL)
//		rq.Header.Set("Content-Type", "application/json")
//		rw := httptest.NewRecorder()
//		b.StartTimer()
//		router.ServeHTTP(rw, rq)
//		res := rw.Result()
//		b.StopTimer()
//		defer res.Body.Close()
//	}
//}
//
//func BenchmarkShortenURLsHandlerJSON(b *testing.B) {
//	type RequestBodyURLs struct {
//		CorrelationID string `json:"correlation_id"`
//		OriginalURL   string `json:"original_url"`
//	}
//
//	cfg := config.LoadConfiguration()
//	storageInstance := storage.Init(cfg.FilePath, cfg.DBPath)
//	middleware.Initialize("debug")
//	router := setupRouter(cfg, storageInstance)
//
//	b.ResetTimer()
//
//	for i := 0; i < b.N; i++ {
//		b.StopTimer() // stop all timers
//		body := []RequestBodyURLs{
//			{
//				CorrelationID: "1",
//				OriginalURL:   "https://ya.com/" + strconv.Itoa(i),
//			},
//			{
//				CorrelationID: "2",
//				OriginalURL:   "https://ya.ru" + strconv.Itoa(i),
//			},
//		}
//		jsonBody, _ := json.Marshal(body)
//		rq := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(string(jsonBody)))
//		rq.Header.Set("Content-Type", "application/json")
//		rw := httptest.NewRecorder()
//		b.StartTimer()
//		router.ServeHTTP(rw, rq)
//		res := rw.Result()
//		b.StopTimer()
//		defer res.Body.Close()
//	}
//}
