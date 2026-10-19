package url_transport_http

import (
	"errors"
	"net/http"

	core_errors "github.com/Sergey-tech9087/golang-url-shortener/internal/core/errors"
	core_logger "github.com/Sergey-tech9087/golang-url-shortener/internal/core/logger"
	core_http_response "github.com/Sergey-tech9087/golang-url-shortener/internal/core/transport/http/response"
)

type GetUrlResponse struct {
	FullUrl string `json:"full_url"`
}

func (h *UrlsHTTPHandler) GetUrl(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	code := r.PathValue("code")

	urlDomain, err := h.urlsService.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			responseHandler.ErrorResponse(err, "url not found")
			return
		}
		responseHandler.ErrorResponse(err, "failed to get url by code")
		return
	}

	response := GetUrlResponse{
		FullUrl: urlDomain.FullUrl,
	}

	responseHandler.JSONResponse(response, http.StatusOK)
}
