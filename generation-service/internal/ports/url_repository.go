package ports

import (
	"context"

	"github.com/uygardeniz/url-shortening/generation-service/internal/types"
)

type URLRepository interface {
	StoreURL(ctx context.Context, url *types.URL) error
}
