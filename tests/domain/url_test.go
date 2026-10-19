package domain_test

import (
	"testing"

	"github.com/Sergey-tech9087/golang-url-shortener/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUrl_WithCode(t *testing.T) {
	original := domain.NewUrlUninitialized("https://google.com")
	updated := original.WithCode("abc123")

	assert.Equal(t, "abc123", updated.Code)
	assert.Equal(t, "https://google.com", updated.FullUrl)

	assert.Empty(t, original.Code)
}

func TestUrl_ShortURL(t *testing.T) {
	tests := []struct {
		name        string
		url         domain.Url
		expected    string
		expectError bool
	}{
		{
			name:        "Valid URL",
			url:         domain.NewUrl(1, "https://go.dev/doc/", "KFd_VPasQl"),
			expected:    "https://go.dev/KFd_VPasQl",
			expectError: false,
		},
		{
			name:        "URL without scheme",
			url:         domain.NewUrl(1, "go.dev/doc", "abc"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.url.ShortURL()

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
