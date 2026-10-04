package changelog

import (
	"os/exec"
	"strings"
)

type GitService struct {
	workingDir string
}

func NewGitService(workingDir string) GitService {
	return GitService{
		workingDir: workingDir,
	}
}

func (g GitService) CurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = g.workingDir

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

// AuthorName returns the trimmed Git user.name value, which may be empty.
// It returns an error if user.name is unset or the Git command fails.
func (g GitService) AuthorName() (string, error) {
	cmd := exec.Command("git", "config", "user.name")
	cmd.Dir = g.workingDir

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}
