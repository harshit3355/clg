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
	for _, src := range sources {
		if src == nil {
			continue
		}

		author, err := src.GetAuthor()
		if err != nil {
			return Author{}, fmt.Errorf("author source %T: %w", src, err)
		}

		if author.IsZero() {
			continue
		}

		return author, nil
	}

	return Author{}, ErrMissingAuthorInfo
}

func (a Author) IsZero() bool {
	return a == (Author{})
}

func (a *Author) normalize() {
	a.Name = strings.TrimSpace(a.Name)
	a.URL = strings.TrimSpace(a.URL)
}

func (a Author) validate() error {
	if a.IsZero() {
		return nil
	}

	var errs []error

	if a.Name == "" {
		errs = append(errs, ErrAuthorNameMissing)
	}

	if a.URL == "" {
		return errors.Join(errs...)
	}

	if err := validation.ValidateURL(a.URL); err != nil {
		errs = append(errs, ErrAuthorURLInvalid)
	}

	return errors.Join(errs...)
}
