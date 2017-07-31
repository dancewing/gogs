package models

import (
	"fmt"
	"io/ioutil"

	git "github.com/gogits/git-module"
	"github.com/gogits/gogs/pkg/pipeline"
)

const JENKINS_CI_FILE = ".jenkins-ci.yml"

type Pipeline struct {
	ID           int64
	OwnerID      int64
	Owner        *User `xorm:"-"`
	Status       int
	RepositoryID int64
	Repository   *Repository `xorm:"-"`
	Jobs         []*Job      `xorm:"-"`
}

type Job struct {
	ID              int64
	JobName         string
	Branch          string
	Commit          string
	PipelineID      int64
	Pipeline        *Pipeline `xorm:"-"`
	Status          int
	Stage           string
	Environment     string
	HookTaskID      int64
	JenkinsHookTask *PipelineHookTask `xorm:"-"`
}

func GetCIFileFromGit(owner *User, repository *Repository) (*pipeline.JobDefinition, error) {
	repoPath := RepoPath(owner.Name, repository.Name)
	repo, err := git.OpenRepository(repoPath)

	if err != nil {
		return nil, err
	}
	commit, err := repo.GetBranchCommit(repository.DefaultBranch)
	if err != nil {
		return nil, err
	}
	treeEntry, err := commit.GetTreeEntryByPath(JENKINS_CI_FILE)
	if err != nil {
		return nil, err
	}
	reader, err := treeEntry.Blob().Data()
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return pipeline.Parse([]byte(data))
}

func CreatePipeline(e Engine, def *pipeline.Pipeline, repo *Repository, hookTaskID int64) error {

	pipeline := Pipeline{
		Status:       0,
		RepositoryID: repo.ID,
	}
	_, err := e.Insert(&pipeline)
	if err != nil {
		return fmt.Errorf("Error GetCIFileFromGit: %v", err)
	}

	if len(def.Jobs) > 0 {
		j := def.Jobs[0]
		job := Job{
			JobName:     j.JobName,
			PipelineID:  pipeline.ID,
			Status:      0,
			Stage:       j.Stage,
			Environment: j.Environment,
			HookTaskID:  hookTaskID,
		}
		_, err = e.Insert(&job)
		if err != nil {
			return fmt.Errorf("Error GetCIFileFromGit: %v", err)
		}
	}

	return nil

}
