package urls_postgres_repository

import (
	core_postgres_pool "github.com/Sergey-tech9087/golang-url-shortener/internal/core/repository/postgres/pool"
)

type UrlsRepository struct {
	pool core_postgres_pool.Pool
}

func NewUrlsRepository(
	pool core_postgres_pool.Pool,
) *UrlsRepository {
	return &UrlsRepository{
		pool: pool,
	}
}
