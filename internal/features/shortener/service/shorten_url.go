package urls_service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"regexp"

	"github.com/Sergey-tech9087/golang-url-shortener/internal/core/domain"
	core_errors "github.com/Sergey-tech9087/golang-url-shortener/internal/core/errors"
)

const (
	alphabet   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	codeLength = 10
	maxRetries = 5
)

var shortCodeRe = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func (s *UrlsService) ShortenUrl(
	ctx context.Context,
	url domain.Url,
) (domain.Url, error) {

	for range maxRetries {
		code, err := generateCode(codeLength)

		if err := Validate(code); err != nil {
			return domain.Url{}, fmt.Errorf("generate code: %w", err)
		}
		urlToSave := url.WithCode(code)
		savedUrl, err := s.urlsRepository.ShortenUrl(ctx, urlToSave)
		if err != nil {
			if errors.Is(err, core_errors.ErrConflict) {
				continue
			}
			return domain.Url{}, fmt.Errorf("create short url: %w", err)
		}
		return savedUrl, nil
	}

	return domain.Url{}, fmt.Errorf("failed to generate unique code after %d attempts: %w", maxRetries, core_errors.ErrNotFound)
}

func generateCode(n int) (string, error) {
	code := make([]byte, n)
	for i := range code {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		code[i] = alphabet[idx.Int64()]
	}
	return string(code), nil
}

func Validate(code string) error {
	fullUrlLen := len([]rune(code))

	if fullUrlLen != 10 {
		return fmt.Errorf(
			"Invalid 'ShortUrl' len: %d; %w",
			fullUrlLen,
			core_errors.ErrInvalidArgument,
		)
	}

	if !shortCodeRe.MatchString(code) {
		return fmt.Errorf(
			"Invalid 'ShortUrl' format: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}
