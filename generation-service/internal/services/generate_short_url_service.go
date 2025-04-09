package services

import (
	"context"
	"crypto/rand"
	"errors"
	"log"
	"math/big"
	"time"

	appErrors "github.com/uygardeniz/url-shortening/generation-service/internal/errors"
	"github.com/uygardeniz/url-shortening/generation-service/internal/ports"
	"github.com/uygardeniz/url-shortening/generation-service/internal/types"
)

type GenerateShortURLService struct {
	urlRepository ports.URLRepository
}

func NewGenerateShortURLService(urlRepository ports.URLRepository) *GenerateShortURLService {
	return &GenerateShortURLService{
		urlRepository: urlRepository,
	}
}

const maxRetries = 5

func (s *GenerateShortURLService) GenerateShortCode(ctx context.Context, originalURL string) (string, error) {
	var shortCode string
	var err error

	// Retry logic for generating a unique short code
	for attempt := 1; attempt <= maxRetries; attempt++ {
		shortCode, err = generateShortCode()
		if err != nil {
			log.Printf("Failed to generate short code: %v", err)
			return "", err
		}

		urlMapping := &types.URL{
			ShortCode:   shortCode,
			OriginalURL: originalURL,
			CreatedAt:   time.Now().Unix(),
		}

		err = s.urlRepository.StoreURL(ctx, urlMapping)

		// Success - no error
		if err == nil {
			if attempt > 1 {
				log.Printf("Successfully stored URL after %d attempts with code: %s", attempt, shortCode)
			}
			return shortCode, nil
		}

		// Handle duplicate short code error
		if errors.Is(err, appErrors.ErrDuplicateShortCode) {
			log.Printf("Attempt %d/%d: Short code collision detected for %s, retrying...",
				attempt, maxRetries, shortCode)
			continue
		}

		log.Printf("Error storing URL: %v", err)
		return "", err
	}

	log.Printf("Failed to generate a unique short code after %d attempts", maxRetries)
	return "", appErrors.ErrDuplicateShortCode
}

func generateShortCode() (string, error) {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	shortCode := make([]byte, 6)

	for i := range shortCode {
		randomIdx, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", appErrors.ErrShortCodeGeneration
		}

		shortCode[i] = chars[randomIdx.Int64()]
	}

	return string(shortCode), nil
}
