package main

import (
	"github.com/LilLebowski/shortener/internal/handlers"
)

func main() {
	router := handlers.SetupRouter()
	err := router.Run(`:8080`)
	if err != nil {
		panic(err)
	}
}
