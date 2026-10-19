package service_test

import (
	"context"
	"testing"

	"github.com/Sergey-tech9087/golang-url-shortener/internal/core/domain"
	core_errors "github.com/Sergey-tech9087/golang-url-shortener/internal/core/errors"
	urls_service "github.com/Sergey-tech9087/golang-url-shortener/internal/features/shortener/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) ShortenUrl(ctx context.Context, url domain.Url) (domain.Url, error) {
	args := m.Called(ctx, url)
	return args.Get(0).(domain.Url), args.Error(1)
}

func (m *MockRepository) GetByCode(ctx context.Context, code string) (domain.Url, error) {
	args := m.Called(ctx, code)
	return args.Get(0).(domain.Url), args.Error(1)
}

func TestUrlsService_ShortenUrl(t *testing.T) {
	ctx := context.Background()
	inputUrl := domain.NewUrlUninitialized("https://google.com")

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := urls_service.NewUrlsService(mockRepo)

		mockRepo.On("ShortenUrl", ctx, mock.AnythingOfType("domain.Url")).
			Return(domain.NewUrl(1, "https://google.com", "mockedCode"), nil).Once()

		result, err := svc.ShortenUrl(ctx, inputUrl)
		require.NoError(t, err)
		assert.Equal(t, "mockedCode", result.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Retry on Collision", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := urls_service.NewUrlsService(mockRepo)

		mockRepo.On("ShortenUrl", ctx, mock.Anything).
			Return(domain.Url{}, core_errors.ErrConflict).Once()
		mockRepo.On("ShortenUrl", ctx, mock.Anything).
			Return(domain.NewUrl(1, "https://google.com", "newCode"), nil).Once()

		result, err := svc.ShortenUrl(ctx, inputUrl)
		require.NoError(t, err)
		assert.Equal(t, "newCode", result.Code)
		mockRepo.AssertExpectations(t)
	})
}
