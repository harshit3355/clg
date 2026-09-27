package cmd

import (
	"github.com/hettiger/clg/internal/changelog"
	"github.com/hettiger/clg/internal/config"
)

type AuthorSourceConfig struct {
	changelog.Author
}

func NewAuthorSourceConfig(c config.Config) AuthorSourceConfig {
	return AuthorSourceConfig{
		changelog.Author{
			Name: c.Author.Name,
			URL:  c.Author.URL,
		},
	}
}

func (s AuthorSourceConfig) GetAuthor() (changelog.Author, error) {
	return changelog.NewAuthor(s.Name, s.URL)
}

type AuthorSourceFlags struct {
	changelog.Author
}

func NewAuthorSourceFlags(name, url string) AuthorSourceFlags {
	return AuthorSourceFlags{
		changelog.Author{
			Name: name,
			URL:  url,
		},
	}
}

func (s AuthorSourceFlags) GetAuthor() (changelog.Author, error) {
	return changelog.NewAuthor(s.Name, s.URL)
}
