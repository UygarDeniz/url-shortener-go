package repositories

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/uygardeniz/url-shortening/reading-service/internal/ports"
	"github.com/uygardeniz/url-shortening/reading-service/internal/types"
)

type CachedRepository struct {
	repository  ports.URLRepository
	redisClient *redis.Client
	cacheTTL    time.Duration
}

func NewCachedRepository(repository ports.URLRepository, redisEndpoint string, cacheTTL time.Duration) *CachedRepository {
	client := redis.NewClient(&redis.Options{
		Addr:     redisEndpoint,
		Password: "",
		DB:       0,
	})

	return &CachedRepository{
		repository:  repository,
		redisClient: client,
		cacheTTL:    cacheTTL,
	}
}

func (r *CachedRepository) GetOriginalURL(ctx context.Context, shortCode string) (*types.URL, error) {
	if shortCode == "" {
		return nil, fmt.Errorf("shortCode cannot be empty")
	}

	cacheKey := "url:" + shortCode
	cachedURL, err := r.redisClient.Get(ctx, cacheKey).Result()

	if err == nil {
		log.Printf("Cache HIT for %s", shortCode)
		return &types.URL{
			ShortCode:   shortCode,
			OriginalURL: cachedURL,
		}, nil
	}

	url, err := r.repository.GetOriginalURL(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	err = r.redisClient.Set(ctx, cacheKey, url.OriginalURL, r.cacheTTL).Err()

	if err != nil {
		log.Printf("Failed to cache %s: %v", shortCode, err)
	}

	return url, nil
}
