package ports

import (
	"context"

	"github.com/uygardeniz/url-shortening/reading-service/internal/types"
)

type URLRepository interface {
	GetOriginalURL(ctx context.Context, shortCode string) (*types.URL, error)
}
