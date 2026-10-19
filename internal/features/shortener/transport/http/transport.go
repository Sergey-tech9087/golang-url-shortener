package url_transport_http

import (
	"context"
	"net/http"

	"github.com/Sergey-tech9087/golang-url-shortener/internal/core/domain"
	core_http_server "github.com/Sergey-tech9087/golang-url-shortener/internal/core/transport/http/server"
)

type UrlsHTTPHandler struct {
	urlsService UrlsService
}

type UrlsService interface {
	ShortenUrl(
		ctx context.Context,
		url domain.Url,
	) (domain.Url, error)
	GetByCode(
		ctx context.Context,
		code string,
	) (domain.Url, error)
}

func NewUrlsHTTPHandler(urlsService UrlsService) *UrlsHTTPHandler {
	return &UrlsHTTPHandler{
		urlsService: urlsService,
	}
}

func (h *UrlsHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/shortenurl",
			Handler: h.ShortenUrl,
		},
		{
			Method:  http.MethodGet,
			Path:    "/urls/{code}",
			Handler: h.GetUrl,
		},
	}
}
