package validation

import (
	"errors"
	"net/url"
	"strings"
)

var ErrMalformedURL = errors.New("malformed URL")
var ErrMissingURLScheme = errors.New("missing URL scheme")
var ErrInvalidURLScheme = errors.New("invalid URL scheme (only http and https are supported)")
var ErrMissingURLHostname = errors.New("missing URL host")

func ValidateURL(rawURL string) error {
	if rawURL != strings.TrimSpace(rawURL) {
		return ErrMalformedURL
	}

	parsedURL, err := url.ParseRequestURI(rawURL)

	if err != nil {
		return ErrMalformedURL
	}

	if parsedURL.Scheme == "" {
		return ErrMissingURLScheme
	}

	if err := ValidateIn("URL scheme", strings.ToLower(parsedURL.Scheme), "http", "https"); err != nil {
		return ErrInvalidURLScheme
	}

	if parsedURL.Hostname() == "" {
		return ErrMissingURLHostname
	}

	return nil
}
