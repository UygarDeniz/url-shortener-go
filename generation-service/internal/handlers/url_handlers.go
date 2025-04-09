package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	appErrors "github.com/uygardeniz/url-shortening/generation-service/internal/errors"
	"github.com/uygardeniz/url-shortening/generation-service/internal/services"
	"github.com/uygardeniz/url-shortening/generation-service/internal/types"
	"github.com/uygardeniz/url-shortening/generation-service/pkg"
)

type URLHandler struct {
	urlService *services.GenerateShortURLService
	validate   *validator.Validate
}

func NewURLHandler(urlService *services.GenerateShortURLService) *URLHandler {
	return &URLHandler{
		urlService: urlService,
		validate:   validator.New(),
	}
}

func (h *URLHandler) HandleGenerateURL(w http.ResponseWriter, r *http.Request) {
	var url types.URLRequest

	err := json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		pkg.SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.validate.Struct(url); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		pkg.SendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Validation error: %s", validationErrors))
		return
	}

	shortCode, err := h.urlService.GenerateShortCode(r.Context(), url.OriginalURL)
	if err != nil {
		log.Printf("Error generating short code: %v", err)
		statusCode := http.StatusInternalServerError
		message := "Failed to generate short URL"

		if errors.Is(err, appErrors.ErrDatabaseAccess) {
			statusCode = http.StatusServiceUnavailable
			message = "Service unavailable"
		} else if errors.Is(err, appErrors.ErrDuplicateShortCode) {
			statusCode = http.StatusConflict
			message = "Failed to generate short URL. Please try again."
		}
		pkg.SendErrorResponse(w, statusCode, message)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := types.URLResponse{
		ShortURL: "localhost:8080" + "/" + shortCode,
	}

	json.NewEncoder(w).Encode(response)
}
