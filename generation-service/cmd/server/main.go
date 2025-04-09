package main

import (
	"log"
	"net/http"

	"github.com/uygardeniz/url-shortening/generation-service/internal/handlers"
	"github.com/uygardeniz/url-shortening/generation-service/internal/repositories"
	"github.com/uygardeniz/url-shortening/generation-service/internal/services"
)

func main() {

	urlRepo, err := repositories.NewDynamoDBRepository("url-mappings")
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	urlService := services.NewGenerateShortURLService(urlRepo)
	urlHandler := handlers.NewURLHandler(urlService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /generate", urlHandler.HandleGenerateURL)

	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	log.Println("Starting URL shortening service on :8081")
	log.Fatal(server.ListenAndServe())
}
