package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/otterize/otterize-cli/src/pkg/errors"
	"os"
	"path/filepath"
)

func GetGitRoot(repo *git.Repository) (string, error) {
	wt, err := repo.Worktree()
	if err != nil {
		return "", errors.Wrap(err)
	}
	return wt.Filesystem.Root(), nil
}

func GetGitRepoInformation(workingDir string) (*LocalGitInformation, error) {
	var err error
	if workingDir == "" {
		workingDir = os.Getenv("PWD")
	}

	repo, err := git.PlainOpenWithOptions(workingDir, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return nil, errors.Wrap(err)
	}

	remotes, err := repo.Remotes()
	if err != nil {
		return nil, errors.Wrap(err)
	}

	headRef, err := repo.Head()
	if err != nil {
		return nil, errors.Wrap(err)
	}

	gitRoot, err := GetGitRoot(repo)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	var gitInfo LocalGitInformation
	gitInfo.Commit = headRef.Hash().String()

	relativePath, err := filepath.Rel(gitRoot, workingDir)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	gitInfo.RelativePath = relativePath

	for _, remote := range remotes {
		if remote.Config().Name == "origin" {
			gitInfo.OriginUrl = remote.Config().URLs[0] // Get the first URL
			break
		}
	}

	return &gitInfo, nil
}
