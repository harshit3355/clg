package changelog

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hettiger/clg/internal/validation"
)

type Author struct {
	Name string
	URL  string
}

type AuthorSource interface {
	GetAuthor() (Author, error)
}

var (
	ErrMissingAuthorInfo = errors.New("missing author information")
	ErrAuthorNameMissing = errors.New("author name missing")
	ErrAuthorURLInvalid  = errors.New("author URL invalid")
)

func NewAuthor(name, url string) (Author, error) {
	author := Author{
		Name: name,
		URL:  url,
	}

	author.normalize()

	if err := author.validate(); err != nil {
		return Author{}, err
	}

	return author, nil
}

func ResolveAuthor(sources []AuthorSource) (Author, error) {
	var errs []error

	for _, src := range sources {
		if src == nil {
			continue
		}

		author, err := src.GetAuthor()
		if err != nil {
			errs = append(errs, fmt.Errorf("author source %T: %w", src, err))
			continue
		}

		author, err = NewAuthor(author.Name, author.URL)
		if err != nil {
			errs = append(errs, fmt.Errorf("author source %T: %w", src, err))
			continue
		}

		return author, nil
	}

	if len(errs) > 0 {
		return Author{}, errors.Join(errs...)
	}

	return Author{}, ErrMissingAuthorInfo
}

func (a *Author) normalize() {
	a.Name = strings.TrimSpace(a.Name)
	a.URL = strings.TrimSpace(a.URL)
}

func (a Author) validate() error {
	if a.Name == "" {
		return ErrAuthorNameMissing
	}

	if a.URL == "" {
		return nil
	}

	if err := validation.ValidateURL(a.URL); err != nil {
		return ErrAuthorURLInvalid
	}

	return nil
}
