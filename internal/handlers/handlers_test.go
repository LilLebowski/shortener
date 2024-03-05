package handlers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateShortURLHandler(t *testing.T) {
	urls = make(map[string]string)

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
			fmt.Printf("\n\nTest %v Body %v\n", test.name, test.param)
			router := SetupRouter()
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
			url:   "",
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
			router := SetupRouter()
			if test.urlID == "found" {
				urls[test.urlID] = test.url
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
