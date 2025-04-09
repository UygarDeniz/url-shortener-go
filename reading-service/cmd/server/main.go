package main

import (
	"log"
	"net/http"
	"time"

	"github.com/uygardeniz/url-shortening/reading-service/internal/handlers"
	"github.com/uygardeniz/url-shortening/reading-service/internal/repositories"
	"github.com/uygardeniz/url-shortening/reading-service/internal/services"
)

func main() {

	dynamoRepo, err := repositories.NewDynamoDBRepository("url-mappings")
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	redisEndpoint := "localhost:6379"
	cacheTTL := 24 * time.Hour
	cachedRepo := repositories.NewCachedRepository(dynamoRepo, redisEndpoint, cacheTTL)

	urlService := services.NewURLService(cachedRepo)
	urlHandler := handlers.NewURLHandler(urlService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{shortCode}", urlHandler.HandleGetOriginalURL)
	server := &http.Server{
		Addr:    ":8082",
		Handler: mux,
	}

	log.Println("Starting URL shortening service on :8082")
	log.Fatal(server.ListenAndServe())
}
