package url_transport_http

import (
	"net/http"

	"github.com/Sergey-tech9087/golang-url-shortener/internal/core/domain"
	core_logger "github.com/Sergey-tech9087/golang-url-shortener/internal/core/logger"
	core_http_request "github.com/Sergey-tech9087/golang-url-shortener/internal/core/transport/http/request"
	core_http_response "github.com/Sergey-tech9087/golang-url-shortener/internal/core/transport/http/response"
)

type ShortenUrlRequest struct {
	FullUrl string `json:"full_url" validate:"required"`
}

type ShortenUrlResponse struct {
	ShortUrl string `json:"short_url"`
}

func (h *UrlsHTTPHandler) ShortenUrl(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request ShortenUrlRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request")

		return
	}

	urlDomain := domainFromDTO(request)

	urlDomain, err := h.urlsService.ShortenUrl(ctx, urlDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to create shortener url")
		return
	}

	response, err := dtoFromDomain(urlDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to build short url")
		return
	}

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(dto ShortenUrlRequest) domain.Url {
	return domain.NewUrlUninitialized(dto.FullUrl)
}

func dtoFromDomain(url domain.Url) (ShortenUrlResponse, error) {
	shortURL, err := url.ShortURL()
	if err != nil {
		return ShortenUrlResponse{}, err
	}
	return ShortenUrlResponse{ShortUrl: shortURL}, nil
}
