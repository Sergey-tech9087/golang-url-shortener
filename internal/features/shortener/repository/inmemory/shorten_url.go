package urls_inmemory_repository

import (
	"context"

	"github.com/Sergey-tech9087/golang-url-shortener/internal/core/domain"
	core_errors "github.com/Sergey-tech9087/golang-url-shortener/internal/core/errors"
)

func (r *UrlsRepository) ShortenUrl(ctx context.Context, url domain.Url) (domain.Url, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.byCode[url.Code]; exists {
		return domain.Url{}, core_errors.ErrConflict
	}

	r.byCode[url.Code] = url
	return url, nil
}
