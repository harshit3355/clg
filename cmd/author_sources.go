package cmd

import (
	"github.com/hettiger/clg/internal/changelog"
	"github.com/hettiger/clg/internal/config"
)

type authorSourceFlags struct {
	name string
	url  string
}

func newAuthorSourceFlags(name, url string) authorSourceFlags {
	return authorSourceFlags{
		name: name,
		url:  url,
	}
}

func (s authorSourceFlags) GetAuthor() (changelog.Author, error) {
	return changelog.NewAuthor(s.name, s.url)
}

type authorSourceConfig struct {
	config config.AuthorConfig
}

func newAuthorSourceConfig(config config.AuthorConfig) authorSourceConfig {
	return authorSourceConfig{
		config: config,
	}
}

func (s authorSourceConfig) GetAuthor() (changelog.Author, error) {
	return changelog.NewAuthor(s.config.Name, s.config.URL)
}

type gitAuthorNameReader interface {
	AuthorName() (string, error)
}

type authorSourceGit struct {
	gitService gitAuthorNameReader
}

func newAuthorSourceGit(gitService gitAuthorNameReader) authorSourceGit {
	return authorSourceGit{
		gitService: gitService,
	}
}

func (s authorSourceGit) GetAuthor() (changelog.Author, error) {
	name, err := s.gitService.AuthorName()
	if err != nil {
		// git author lookup is best-effort; ignore errors.
		return changelog.Author{}, nil
	}

	return changelog.NewAuthor(name, "")
}
