package cmd

import (
	"github.com/hettiger/clg/internal/changelog"
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
		return changelog.Author{}, err
	}

	return changelog.NewAuthor(name, "")
}
