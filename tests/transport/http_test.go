package transport_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sergey-tech9087/golang-url-shortener/internal/core/domain"
	core_errors "github.com/Sergey-tech9087/golang-url-shortener/internal/core/errors"
	core_logger "github.com/Sergey-tech9087/golang-url-shortener/internal/core/logger"
	url_transport_http "github.com/Sergey-tech9087/golang-url-shortener/internal/features/shortener/transport/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) ShortenUrl(ctx context.Context, url domain.Url) (domain.Url, error) {
	args := m.Called(ctx, url)
	return args.Get(0).(domain.Url), args.Error(1)
}

func (m *MockService) GetByCode(ctx context.Context, code string) (domain.Url, error) {
	args := m.Called(ctx, code)
	return args.Get(0).(domain.Url), args.Error(1)
}

func prepareRequestWithLogger(t *testing.T, req *http.Request) *http.Request {
	t.Helper()

	logConfig := core_logger.LoggerConfig{
		Level:  "debug",
		Folder: t.TempDir(),
	}

	logger, err := core_logger.NewLogger(logConfig)
	require.NoError(t, err)

	t.Cleanup(func() {
		logger.Close()
	})

	ctx := core_logger.ToContext(req.Context(), logger)
	return req.WithContext(ctx)
}

func TestHTTPHandler_GetUrl(t *testing.T) {
	mockSvc := new(MockService)
	handler := url_transport_http.NewUrlsHTTPHandler(mockSvc)

	t.Run("Not Found", func(t *testing.T) {
		mockSvc.On("GetByCode", mock.Anything, "notfound").Return(domain.Url{}, core_errors.ErrNotFound)

		mux := http.NewServeMux()
		mux.HandleFunc("GET /api/v1/urls/{code}", handler.GetUrl)

		req := httptest.NewRequest("GET", "/api/v1/urls/notfound", nil)
		req = prepareRequestWithLogger(t, req)

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Success", func(t *testing.T) {
		mockSvc.On("GetByCode", mock.Anything, "validCode").Return(domain.NewUrl(1, "https://google.com", "validCode"), nil)

		mux := http.NewServeMux()
		mux.HandleFunc("GET /api/v1/urls/{code}", handler.GetUrl)

		req := httptest.NewRequest("GET", "/api/v1/urls/validCode", nil)
		req = prepareRequestWithLogger(t, req)

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"full_url":"https://google.com"}`, w.Body.String())
	})
}

func TestHTTPHandler_ShortenUrl(t *testing.T) {
	mockSvc := new(MockService)
	handler := url_transport_http.NewUrlsHTTPHandler(mockSvc)

	t.Run("Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/shortenurl", bytes.NewBufferString(`{"full_url": bad json}`))
		req = prepareRequestWithLogger(t, req)

		w := httptest.NewRecorder()
		handler.ShortenUrl(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Success", func(t *testing.T) {
		mockSvc.On("GetByCode", mock.Anything, "validCode").Return(domain.NewUrl(1, "https://google.com", "validCode"), nil)

		mux := http.NewServeMux()
		mux.HandleFunc("GET /api/v1/urls/{code}", handler.GetUrl)

		req := httptest.NewRequest("GET", "/api/v1/urls/validCode", nil)
		req = prepareRequestWithLogger(t, req)

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		assert.JSONEq(t, `{"full_url":"https://google.com"}`, w.Body.String())
	})
}
