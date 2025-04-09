package handlers

import (
	"errors"
	"net/http"

	appErrors "github.com/uygardeniz/url-shortening/reading-service/internal/errors"
	"github.com/uygardeniz/url-shortening/reading-service/internal/services"
	"github.com/uygardeniz/url-shortening/reading-service/pkg"
)

type URLHandler struct {
	urlService *services.URLService
}

func NewURLHandler(urlService *services.URLService) *URLHandler {
	return &URLHandler{
		urlService: urlService,
	}
}

func (h *URLHandler) HandleGetOriginalURL(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("shortCode")
	ctx := r.Context()
	originalURL, err := h.urlService.GetOriginalURL(ctx, shortCode)

	if err != nil {

		switch {
		case errors.Is(err, appErrors.ErrURLNotFound):
			pkg.SendErrorResponse(w, http.StatusNotFound, "URL not found")

		case errors.Is(err, appErrors.ErrEmptyShortCode):
			pkg.SendErrorResponse(w, http.StatusBadRequest, "Short code cannot be empty")
		case errors.Is(err, appErrors.ErrDatabaseAccess),
			errors.Is(err, appErrors.ErrDataFormat),
			errors.Is(err, appErrors.ErrInternalServerError):
			pkg.SendErrorResponse(w, http.StatusInternalServerError, "An error occurred while processing your request")
		default:
			pkg.SendErrorResponse(w, http.StatusInternalServerError, "An unexpected error occurred")
		}
		return

	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}
