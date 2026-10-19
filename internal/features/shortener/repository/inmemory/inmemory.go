package urls_inmemory_repository

import (
	"sync"

	"github.com/Sergey-tech9087/golang-url-shortener/internal/core/domain"
)

type UrlsRepository struct {
	mu     sync.RWMutex
	byCode map[string]domain.Url
}

func NewUrlsRepository() *UrlsRepository {
	return &UrlsRepository{
		byCode: make(map[string]domain.Url),
	}
}
