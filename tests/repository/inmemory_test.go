package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Sergey-tech9087/golang-url-shortener/internal/core/domain"
	core_errors "github.com/Sergey-tech9087/golang-url-shortener/internal/core/errors"
	urls_inmemory_repository "github.com/Sergey-tech9087/golang-url-shortener/internal/features/shortener/repository/inmemory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryRepository(t *testing.T) {
	ctx := context.Background()
	repo := urls_inmemory_repository.NewUrlsRepository()

	url := domain.NewUrl(0, "https://example.com", "testCode")

	savedUrl, err := repo.ShortenUrl(ctx, url)
	require.NoError(t, err)
	assert.Equal(t, "testCode", savedUrl.Code)

	_, err = repo.ShortenUrl(ctx, url)
	require.Error(t, err)
	assert.True(t, errors.Is(err, core_errors.ErrConflict), "должна вернуться ошибка конфликта")

	foundUrl, err := repo.GetByCode(ctx, "testCode")
	require.NoError(t, err)
	assert.Equal(t, "https://example.com", foundUrl.FullUrl)

	_, err = repo.GetByCode(ctx, "nonExistent")
	require.Error(t, err)
	assert.True(t, errors.Is(err, core_errors.ErrNotFound))
}
