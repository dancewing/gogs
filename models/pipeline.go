package models

import (
	"fmt"
	"io/ioutil"

	"time"

	"strings"

	"github.com/go-xorm/xorm"
	git "github.com/gogits/git-module"
	api "github.com/gogits/go-gogs-client"
	"github.com/gogits/gogs/pkg/jenkins_ci_parser"
	"github.com/gogits/gogs/pkg/jenkins_client"
	"github.com/kataras/iris/core/errors"
	log "gopkg.in/clog.v1"
)

const JENKINS_CI_YAML = ".jenkins-ci.yml"
const JENKINS_CI_JSON = ".jenkins-ci.json"

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
	Status       string
	RepositoryID int64
	Repository   *Repository `xorm:"-"`
	Jobs         []*Job      `xorm:"-"`

	JenkinsJobName string
	DeliveryUUID   string

	ExecuteLog string

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
	ID int64

	PipelineID   int64
	Pipeline     *Pipeline `xorm:"-"`
	RepositoryID int64

	Status      string
	Stage       string
	Environment string

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

func CountJobByStatus(pipelineID int64, status string) int64 {
	sess := x.Where("id > 0")

	if pipelineID > 0 {
		sess.And("pipeline_id = ?", pipelineID)
	}

	if status != "" {
		sess.And("status = ?", status)
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

func getPipelineByDeliveryID(e Engine, deliveryID string) (*Pipeline, error) {
	m := &Pipeline{
		DeliveryUUID: deliveryID,
	}
	has, err := e.Get(m)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.New(fmt.Sprintf("Pipeline not exist with deliveryID- %s", deliveryID))
	}
	return m, nil
}

func GetPipeline(pipelineID int64) (*Pipeline, error) {
	return getPipelineByID(x, pipelineID)
}

func GetPipelineByDeliveryID(deliveryID string) (*Pipeline, error) {
	return getPipelineByDeliveryID(x, deliveryID)
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

func GetJobByStep(step string, pipelineID int64) (*Job, error) {
	return getJobByStep(x, step, pipelineID)
}

func getJobByStep(e Engine, step string, pipelineID int64) (*Job, error) {
	m := &Job{
		PipelineID: pipelineID,
		Stage:      step,
	}
	has, err := e.Get(m)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.New(fmt.Sprintf("Job not exist with pipelineID - %d, stage : %s ", pipelineID, step))
	}
	return m, nil
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

func UpdateJob(job *Job) error {
	_, err := x.Id(job.ID).AllCols().Update(job)
	return err
}

func GetCIFileFromGit(repository *Repository, branch string, fileType string) (string, error) {
	repoPath := RepoPath(repository.Owner.Name, repository.Name)
	repo, err := git.OpenRepository(repoPath)

	if err != nil {
		return "", err
	}
	var br string
	if branch == "" {
		br = repository.DefaultBranch
	} else {
		br = branch
	}
	commit, err := repo.GetBranchCommit(br)
	if err != nil {
		return "", err
	}

	var fileName string
	if fileType == "yaml" {
		fileName = JENKINS_CI_YAML
	} else {
		fileName = JENKINS_CI_JSON
	}

	treeEntry, err := commit.GetTreeEntryByPath(fileName)
	if err != nil {
		return "", err
	}
	reader, err := treeEntry.Blob().Data()
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func CreatePipeline(pipeline *Pipeline) (*Pipeline, error) {

	_, err := x.Insert(pipeline)
	if err != nil {
		return nil, fmt.Errorf("Error CreatePipeline: %v", err)
	}
	return pipeline, nil
}

func CreateJob(job *Job) (*Job, error) {

	_, err := x.Insert(job)
	if err != nil {
		return nil, fmt.Errorf("Error CreateJob: %v", err)
	}
	return job, nil
}

func UpdatePipelineStatus(pipeline *Pipeline) (*Pipeline, error) {

	all := CountJobByStatus(pipeline.ID, "")
	success := CountJobByStatus(pipeline.ID, "success")
	failed := CountJobByStatus(pipeline.ID, "failed")
	canceled := CountJobByStatus(pipeline.ID, "canceled")

	if all == success {
		pipeline.Status = "success"
	}
	if failed >= 1 {
		pipeline.Status = "failed"
	}

	if canceled >= 1 {
		pipeline.Status = "canceled"
	}

	_, err := x.ID(pipeline.ID).Cols("status").Update(pipeline)
	if err != nil {
		return nil, fmt.Errorf("Error UpdatePipeline: %v", err)
	}
	return pipeline, nil
}

func PreviewPipelineScript(repository *Repository, branch string) (string, error) {
	config, err := GetConfiguration(JENKINS, repository.ID)

	if err != nil {
		return "", err
	}

	if !config.IsActive {
		return "", errors.New("Jenkins Config is not active")
	}

	jenkinsCfg := config.ToJenkinsServiceConfigEdit()

	var content string
	if jenkinsCfg.JenkinsFileType == "yaml" {
		content, err = GetCIFileFromGit(repository, branch, "yaml")

		if err != nil {
			return "", err
		}

		parser := jenkins_ci_parser.NewGitLabCiYamlParser([]byte(content))

		ci, err := parser.ParseYaml()

		if err != nil {
			return "", err
		}

		writer := jenkins_ci_parser.NewPipelineWriter(true)

		p := ci.Pipeline.FilterStages(branch, ci.DefaultEnvironment().Name)

		p.Writer(writer)

		return writer.String(), nil

	} else {

	}

	return "", nil

}

func RunPipeline(repository *Repository, branch string, environment string, payload api.Payloader) error {
	config, err := GetConfiguration(JENKINS, repository.ID)

	if err != nil {
		return err
	}

	if !config.IsActive {
		return errors.New("Jenkins Config is not active")
	}

	jenkinsCfg := config.ToJenkinsServiceConfigEdit()

	var content string
	if jenkinsCfg.JenkinsFileType == "yaml" {
		content, err = GetCIFileFromGit(repository, branch, "yaml")

		if err != nil {
			return err
		}

		parser := jenkins_ci_parser.NewGitLabCiYamlParser([]byte(content))

		ci, err := parser.ParseYaml()

		if err != nil {
			return err
		}

		writer := jenkins_ci_parser.NewPipelineWriter(true)

		if environment == "" {
			environment = ci.DefaultEnvironment().Name
		}

		p := ci.Pipeline.FilterStages(branch, environment)

		p.Writer(writer)

		//Generate Job
		jobName := generateJobName(repository, branch, environment)

		CheckJenkinsJob(jenkinsCfg.JenkinsHost, jenkinsCfg.JenkinsUser, jenkinsCfg.JenkinsToken, jobName, writer.String())

		//Check Remote Job by Jenkins Client
		task, err := runServiceTask(repository, payload, config, jobName)

		if err != nil {
			return err
		}

		pipeline := &Pipeline{
			Status:         "running",
			RepositoryID:   repository.ID,
			DeliveryUUID:   task.UUID,
			JenkinsJobName: jobName,
		}

		_, err = CreatePipeline(pipeline)

		if err != nil {
			return err
		}

		for _, stage := range p.Stages {
			job := &Job{
				PipelineID:   pipeline.ID,
				Stage:        stage.Name,
				RepositoryID: pipeline.RepositoryID,
			}
			_, err = CreateJob(job)
			if err != nil {
				return err
			}
		}

	} else {

	}

	return nil
}
func generateJobName(repository *Repository, branch string, environment string) string {
	return strings.Join([]string{repository.Owner.Name, repository.Name, branch, environment}, "_")
}

func GetLastCommit(repository *Repository, branch string) (*git.Commit, error) {
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
	return commit, nil
}

func CheckJenkinsJob(host string, apiUser string, apiToken string, jobName string, scripts string) error {

	jenkins := jenkins_client.NewJenkins(&jenkins_client.Auth{Username: apiUser, ApiToken: apiToken}, host)

	_, err := jenkins.GetJob(jobName)

	if err != nil {
		//create Job

		template, err := jenkins_client.NewWorkflowJobTemplate(jenkins)

		if err != nil {
			return err
		}

		template.Definition.Script = scripts

		if err != nil {
			jenkins.UpdateJob(template, jobName)
		} else {
			jenkins.CreateJob(template, jobName)
		}

	}

	return nil
}
