package models

import (
	"fmt"
	"io/ioutil"

	"time"

	"github.com/go-xorm/xorm"
	git "github.com/gogits/git-module"
	"github.com/gogits/gogs/pkg/pipeline"
	log "gopkg.in/clog.v1"
)

const JENKINS_CI_FILE = ".jenkins-ci.yml"

type ErrPipelineNotExist struct {
	ID int64
}

func IsErrPipelineNotExist(err error) bool {
	_, ok := err.(ErrPipelineNotExist)
	return ok
}

func (err ErrPipelineNotExist) Error() string {
	return fmt.Sprintf("pipeline does not exist [id: %d ]", err.ID)
}

type ErrJobNotExist struct {
	ID     int64
	TaskID int64
}

func IsErrJobNotExist(err error) bool {
	_, ok := err.(ErrPipelineNotExist)
	return ok
}

func (err ErrJobNotExist) Error() string {
	return fmt.Sprintf("job does not exist [id: %d, taskID: %d ]", err.ID, err.TaskID)
}

type Pipeline struct {
	ID           int64
	OwnerID      int64
	Owner        *User `xorm:"-"`
	Status       int
	RepositoryID int64
	Repository   *Repository `xorm:"-"`
	Jobs         []*Job      `xorm:"-"`

	Created     time.Time `xorm:"-"`
	CreatedUnix int64
	Updated     time.Time `xorm:"-"`
	UpdatedUnix int64
}

func (p *Pipeline) BeforeInsert() {
	p.CreatedUnix = time.Now().Unix()
	p.UpdatedUnix = p.CreatedUnix
}

func (p *Pipeline) BeforeUpdate() {
	p.UpdatedUnix = time.Now().Unix()
}

func (p *Pipeline) loadAttributes(e Engine) (err error) {

	jobs, err := ListJobs(&JobOptions{
		PipelineID:   p.ID,
		RepositoryID: p.RepositoryID,
		Page:         1,
		PageSize:     10,
	})

	p.Jobs = jobs

	return err
}

func (p *Pipeline) LoadAttributes() error {
	return p.loadAttributes(x)
}

type PipelineOptions struct {
	RepositoryID int64
	Page         int
	PageSize     int
}

type JobOptions struct {
	RepositoryID int64
	PipelineID   int64
	Page         int
	PageSize     int
}

type Job struct {
	ID             int64
	JobName        string
	NextJobName    string
	JenkinsJobName string
	Branch         string
	Commit         string
	PipelineID     int64
	Pipeline       *Pipeline `xorm:"-"`
	RepositoryID   int64
	Repository     *Repository `xorm:"-"`
	Status         int
	Stage          string
	Environment    string
	DeliveryUUID   string
	HookTaskID     int64
	HookTask       *PipelineHookTask `xorm:"-"`

	JobURL string `xorm:"TEXT"`

	Created     time.Time `xorm:"-"`
	CreatedUnix int64
	Updated     time.Time `xorm:"-"`
	UpdatedUnix int64
}

type JobCallBack struct {
	ID             int64
	DeliveryUUID   string
	JobID          int64
	Job            *Job `xorm:"-"`
	IsSucceed      bool
	RequestContent string `xorm:"TEXT"`
}

func (p *Job) BeforeInsert() {
	p.CreatedUnix = time.Now().Unix()
	p.UpdatedUnix = p.CreatedUnix
}

func (p *Job) BeforeUpdate() {
	p.UpdatedUnix = time.Now().Unix()
}

func (p *Job) loadAttributes(e Engine) (err error) {

	if p.RepositoryID > 0 {
		p.Repository, err = GetRepositoryByID(p.RepositoryID)
	}

	if p.HookTaskID > 0 {
		p.HookTask, err = GetPipelineHookTask(p.HookTaskID)
	}

	return nil
}

func (p *Job) LoadAttributes() error {
	return p.loadAttributes(x)
}

func (p *Job) FindNextJob() (*Job, error) {
	m := &Job{
		PipelineID: p.PipelineID,
		JobName:    p.NextJobName,
	}
	has, err := x.Get(m)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrJobNotExist{ID: p.ID, TaskID: 0}
	}
	return m, nil
}

func CountPipeline(repositoryID int64) int64 {
	sess := x.Where("id > 0")

	if repositoryID > 0 {
		sess.And("repository_id = ?", repositoryID)
	}

	count, err := sess.Count(new(Pipeline))
	if err != nil {
		log.Error(4, "CountPipeline: %v", err)
	}
	return count
}

func CountJob(repositoryID int64) int64 {
	sess := x.Where("id > 0")

	if repositoryID > 0 {
		sess.And("repository_id = ?", repositoryID)
	}

	count, err := sess.Count(new(Job))
	if err != nil {
		log.Error(4, "CountJob: %v", err)
	}
	return count
}

// GetUserRepositories returns a list of repositories of given user.
func ListPipelines(opts *PipelineOptions) ([]*Pipeline, error) {
	sess := x.Where("repository_id=?", opts.RepositoryID).Desc("created_unix")

	if opts.Page <= 0 {
		opts.Page = 1
	}
	sess.Limit(opts.PageSize, (opts.Page-1)*opts.PageSize)

	pipelines := make([]*Pipeline, 0, opts.PageSize)
	return pipelines, sess.Find(&pipelines)
}

func getPipelineByID(e Engine, pipelineID int64) (*Pipeline, error) {
	m := &Pipeline{
		ID: pipelineID,
	}
	has, err := e.Get(m)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrPipelineNotExist{pipelineID}
	}
	return m, nil
}

func GetPipeline(pipelineID int64) (*Pipeline, error) {
	return getPipelineByID(x, pipelineID)
}

func ListJobs(opts *JobOptions) ([]*Job, error) {

	var sess *xorm.Session
	if opts.PipelineID > 0 {
		sess = x.Where("repository_id=? and pipeline_id = ? ", opts.RepositoryID, opts.PipelineID).Desc("id")
	} else {
		sess = x.Where("repository_id=?", opts.RepositoryID).Desc("id")
	}

	if opts.Page <= 0 {
		opts.Page = 1
	}
	sess.Limit(opts.PageSize, (opts.Page-1)*opts.PageSize)

	jobs := make([]*Job, 0, opts.PageSize)
	return jobs, sess.Find(&jobs)
}

func GetJob(jobID int64) (*Job, error) {
	return getJobByID(x, jobID)
}

func getJobByID(e Engine, jobID int64) (*Job, error) {
	m := &Job{
		ID: jobID,
	}
	has, err := e.Get(m)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrJobNotExist{ID: jobID, TaskID: 0}
	}
	return m, nil
}

func UpdateJob(job *Job) (error) {
	_, err := x.Id(job.ID).AllCols().Update(job)
	return err
}

func GetJobByTask(taskID int64) (*Job, error) {
	m := &Job{
		HookTaskID: taskID,
	}
	has, err := x.Get(m)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrJobNotExist{ID: 0, TaskID: taskID}
	}
	return m, nil
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

func CreatePipeline(e Engine, def *pipeline.JobDefinition, repo *Repository) (*Job, error) {

	pipeline := Pipeline{
		Status:       0,
		RepositoryID: repo.ID,
	}
	_, err := e.Insert(&pipeline)
	if err != nil {
		return nil, fmt.Errorf("Error GetCIFileFromGit: %v", err)
	}

	var firstJob *Job
	for index := 0; index < len(def.Pipeline.Jobs); index++ {
		j := def.Pipeline.Jobs[index]

		job := Job{
			JobName:        j.JobName,
			PipelineID:     pipeline.ID,
			Status:         0,
			Stage:          j.Stage,
			Environment:    j.Environment,
			RepositoryID:   repo.ID,
			NextJobName:    j.NextJobName,
			JenkinsJobName: j.JenkinsJob,
			JobURL:         def.GetJobURL(j),
		}

		_, err = e.Insert(&job)
		if err != nil {
			return nil, fmt.Errorf("Error GetCIFileFromGit: %v", err)
		}
		if index == 0 {
			firstJob = &job
		}
	}

	return firstJob, nil

}
