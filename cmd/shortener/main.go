package main

import (
	"github.com/LilLebowski/shortener/internal/handlers"
	"net/http"
)

var urls map[string]string

func handler(rw http.ResponseWriter, rq *http.Request) {
	switch rq.Method {
	case http.MethodPost:
		handlers.CreateShortURLHandler(rw, rq, urls)
	case http.MethodGet:
		handlers.GetShortURLHandler(rw, rq, urls)
	}
}

func main() {
	urls = make(map[string]string)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
