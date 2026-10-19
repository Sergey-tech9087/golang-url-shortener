package domain

import (
	"fmt"
	"net/url"
	"strings"
)

type Url struct {
	ID      int
	FullUrl string
	Code    string
}

func NewUrl(
	id int,
	fullUrl string,
	code string,
) Url {
	return Url{
		ID:      id,
		FullUrl: fullUrl,
		Code:    code,
	}
}

func NewUrlUninitialized(
	fullUrl string,
) Url {
	return NewUrl(
		UninitializedID,
		fullUrl,
		UninitializedShortUrl,
	)
}

func (u Url) WithCode(code string) Url {
	u.Code = code
	return u
}

func (u Url) ShortURL() (string, error) {
	parsed, err := url.Parse(u.FullUrl)
	if err != nil {
		return "", fmt.Errorf("parse url: %w", err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("url must have scheme and host: %q", u.FullUrl)
	}

	parsed.Path = "/" + strings.TrimSpace(u.Code)
	parsed.RawPath = ""
	parsed.RawQuery = ""
	parsed.Fragment = ""
	parsed.ForceQuery = false

	return parsed.String(), nil
}
