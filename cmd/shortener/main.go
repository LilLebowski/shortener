package main

import (
	"github.com/LilLebowski/shortener/internal/handlers"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	router := handlers.SetupRouter()
	serverAddr := os.Getenv("HOST") + ":" + os.Getenv("PORT")
	routerErr := router.Run(serverAddr)
	if routerErr != nil {
		panic(routerErr)
	}
}
