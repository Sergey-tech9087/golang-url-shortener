package urls_service

import (
	"context"

	"github.com/Sergey-tech9087/golang-url-shortener/internal/core/domain"
)

type UrlsService struct {
	urlsRepository UrlsRepository
}

type UrlsRepository interface {
	ShortenUrl(
		ctx context.Context,
		url domain.Url,
	) (domain.Url, error)
	GetByCode(
		ctx context.Context,
		code string,
	) (domain.Url, error)
}

func NewUrlsService(
	urlsRepository UrlsRepository,
) *UrlsService {
	return &UrlsService{
		urlsRepository: urlsRepository,
	}
}
