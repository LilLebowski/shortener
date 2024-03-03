package handlers

import (
	"fmt"
	"github.com/LilLebowski/shortener/internal/utils"
	"io"
	"net/http"
)

func CreateShortURLHandler(rw http.ResponseWriter, rq *http.Request, urls map[string]string) {
	reqBody, err := io.ReadAll(rq.Body)
	if err != nil {
		fmt.Printf("could not read request body: %s\n", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	reqBodyString := string(reqBody)
	fmt.Printf("request body: %s\n", reqBodyString)
	if reqBodyString != "" {
		res, encodeErr := utils.EncodeURL(reqBodyString)
		if encodeErr == nil {
			urls[res] = reqBodyString
			rw.Header().Set("Content-Type", "text/plain")
			rw.WriteHeader(http.StatusCreated)
			_, writeErr := rw.Write([]byte("http://localhost:8080/" + res))
			if writeErr != nil {
				rw.WriteHeader(http.StatusBadRequest)
			}
		} else {
			rw.WriteHeader(http.StatusBadRequest)
		}
	} else {
		rw.WriteHeader(http.StatusBadRequest)
	}
}

func GetShortURLHandler(rw http.ResponseWriter, rq *http.Request, urls map[string]string) {
	fmt.Printf("current session: %s\n", urls)
	urlID := rq.URL.String()[1:]
	fmt.Printf("url id: %s\n", urlID)
	if value, ok := urls[urlID]; ok {
		fmt.Printf("found url: %s\n", value)
		rw.Header().Set("Location", value)
		rw.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		rw.Header().Set("Location", "Not found")
		rw.WriteHeader(http.StatusBadRequest)
	}
}
