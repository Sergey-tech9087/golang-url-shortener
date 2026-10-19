package urls_postgres_repository

import (
	"context"
	"fmt"

	"github.com/Sergey-tech9087/golang-url-shortener/internal/core/domain"
)

type UrlModel struct {
	ID      int
	FullUrl string
	Code    string
}

func (r *UrlsRepository) ShortenUrl(
	ctx context.Context,
	url domain.Url,
) (domain.Url, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO urlShortenerApp.urls (full_url, code)
	VALUES ($1, $2)
	RETURNING id, full_url, code;
	`

	row := r.pool.QueryRow(ctx, query, url.FullUrl, url.Code)

	var urlModel UrlModel

	err := row.Scan(
		&urlModel.ID,
		&urlModel.FullUrl,
		&urlModel.Code,
	)

	if err != nil {
		return domain.Url{}, fmt.Errorf("Scan error: %w", err)
	}

	urlDomain := domain.NewUrl(
		urlModel.ID,
		urlModel.FullUrl,
		urlModel.Code,
	)

	return urlDomain, nil
}
