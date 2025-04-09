package services

import (
	"context"

	"github.com/uygardeniz/url-shortening/reading-service/internal/ports"
)

type URLService struct {
	urlRepository ports.URLRepository
}

func NewURLService(urlRepository ports.URLRepository) *URLService {
	return &URLService{
		urlRepository: urlRepository,
	}
}

func (s *URLService) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	url, err := s.urlRepository.GetOriginalURL(ctx, shortCode)
	if err != nil {
		return "", err
	}
	return url.OriginalURL, nil
}
