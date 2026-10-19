package urls_service

import (
	"context"
	"fmt"

	"github.com/Sergey-tech9087/golang-url-shortener/internal/core/domain"
)

func (s *UrlsService) GetByCode(
	ctx context.Context,
	code string,
) (domain.Url, error) {
	if err := validate(code); err != nil {
		return domain.Url{}, fmt.Errorf("validate code: %w", err)
	}

	url, err := s.urlsRepository.GetByCode(ctx, code)
	if err != nil {
		return domain.Url{}, fmt.Errorf("get url by code: %w", err)
	}

	return url, nil
}
