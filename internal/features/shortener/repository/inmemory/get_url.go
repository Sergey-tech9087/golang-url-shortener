package urls_inmemory_repository

import (
	"context"

	"github.com/Sergey-tech9087/golang-url-shortener/internal/core/domain"
	core_errors "github.com/Sergey-tech9087/golang-url-shortener/internal/core/errors"
)

func (r *UrlsRepository) GetByCode(ctx context.Context, code string) (domain.Url, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, exists := r.byCode[code]
	if !exists {
		return domain.Url{}, core_errors.ErrNotFound
	}
	return u, nil
}
