package urls_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Sergey-tech9087/golang-url-shortener/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

func (r *UrlsRepository) GetByCode(
	ctx context.Context,
	code string,
) (domain.Url, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, full_url, code
	FROM urlShortenerApp.urls
	WHERE code = $1;
	`

	row := r.pool.QueryRow(ctx, query, code)

	var urlModel UrlModel

	err := row.Scan(
		&urlModel.ID,
		&urlModel.FullUrl,
		&urlModel.Code,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Url{}, pgx.ErrNoRows
		}
		return domain.Url{}, fmt.Errorf("select url by code: %w", err)
	}

	urlDomain := domain.NewUrl(
		urlModel.ID,
		urlModel.FullUrl,
		urlModel.Code,
	)

	return urlDomain, nil
}
