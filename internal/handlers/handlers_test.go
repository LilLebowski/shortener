package handlers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var urls map[string]string

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
			param := strings.NewReader(test.param)
			rq := httptest.NewRequest(http.MethodPost, "/", param)
			rw := httptest.NewRecorder()
			CreateShortURLHandler(rw, rq, urls)

			res := rw.Result()
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					panic(err)
				}
			}(res.Body)
			fmt.Printf("want code = %d StatusCode %d\n", test.want.code, res.StatusCode)
			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}

func TestGetShortURLHandler(t *testing.T) {
	urls = make(map[string]string)

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
			name:  "POST 1. body doesn't consist of data",
			urlID: "notfound",
			url:   "",
			want: want{
				code: 400,
			},
		},
		{
			name:  "POST 2. body consist of data",
			urlID: "found",
			url:   "https://ya.ru",
			want: want{
				code: 307,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.urlID == "found" {
				urls[test.urlID] = test.url
			}
			fmt.Printf("\n\nTest %v urlID %v url %v\n", test.name, test.urlID, test.url)
			rq := httptest.NewRequest(http.MethodGet, "/"+test.urlID, nil)
			rw := httptest.NewRecorder()
			GetShortURLHandler(rw, rq, urls)

			res := rw.Result()
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					panic(err)
				}
			}(res.Body)
			fmt.Printf("want code = %d StatusCode %d\n", test.want.code, res.StatusCode)
			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}
