package git

import (
	"github.com/go-git/go-git/v5"
	"os"
	"path/filepath"
)

func GetGitRoot(repo *git.Repository) (string, error) {
	wt, err := repo.Worktree()
	if err != nil {
		return "", err
	}
	return wt.Filesystem.Root(), nil
}

func GetGitRepoInformation(workingDir string) (*LocalGitInformation, error) {
	if workingDir == "" {
		workingDir = os.Getenv("PWD")
	}

	repo, err := git.PlainOpenWithOptions(workingDir, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return nil, err
	}

	remotes, err := repo.Remotes()
	if err != nil {
		return nil, err
	}

	headRef, err := repo.Head()
	if err != nil {
		return nil, err
	}

	gitRoot, err := GetGitRoot(repo)
	if err != nil {
		return nil, err
	}

	var gitInfo LocalGitInformation
	gitInfo.Commit = headRef.Hash().String()

	relativePath, err := filepath.Rel(gitRoot, workingDir)
	if err != nil {
		return nil, err
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
