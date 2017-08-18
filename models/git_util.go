package models

import (
	"io/ioutil"

	git "github.com/gogits/git-module"
)

func GetFileFromGit(repository *Repository, branch string, fileName string) ([]byte, error) {
	repoPath := RepoPath(repository.Owner.Name, repository.Name)
	repo, err := git.OpenRepository(repoPath)

	if err != nil {
		return nil, err
	}
	var br string
	if branch == "" {
		br = repository.DefaultBranch
	} else {
		br = branch
	}
	commit, err := repo.GetBranchCommit(br)
	if err != nil {
		return nil, err
	}

	treeEntry, err := commit.GetTreeEntryByPath(fileName)
	if err != nil {
		return nil, err
	}
	reader, err := treeEntry.Blob().Data()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(reader)
}
