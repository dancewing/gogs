package repo

import (
	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/pkg/context"
)

const (
	PIPELINE_CREATE = "repo/pipeline/new_pipeline"
	PIPELINES       = "repo/pipeline/list"
	PIPELINE_JOBS   = "repo/pipeline/jobs"
	PIPELINE_ENV    = "repo/pipeline/environments"
	PIPELINE_NONE   = "repo/pipeline/no_pipeline"
)

func RetrievePipeline(c *context.Context) {

	c.Data["PageIsPipelines"] = true

	jobDef, err := models.GetCIFileFromGit(c.Repo.Owner, c.Repo.Repository)
	if err != nil {
		c.Handle(200, "RetrievePipeline.GetCIFileFromGit", err)
		return
	}

	if jobDef == nil {
		c.HTML(200, PIPELINE_NONE)
		return
	}

	c.Data["JobDef"] = jobDef
	c.Data["PageIsPipelines"] = true
}

func ListPipelines(c *context.Context) {

	c.Data["PageIsPipelineList"] = true

	c.HTML(200, PIPELINES)

}

func NewPipeline(c *context.Context) {

	c.Data["PageIsPipelineList"] = true

	c.HTML(200, PIPELINE_CREATE)

}

func NewPipelinePost(c *context.Context) {

	c.Data["PageIsPipelineList"] = true

	c.HTML(200, PIPELINES)

}

func ListJobs(c *context.Context) {
	c.Data["PageIsJobList"] = true

	c.HTML(200, PIPELINE_JOBS)
}

func ListEnvironments(c *context.Context) {

	c.Data["PageIsEnvironmentList"] = true

	c.HTML(200, PIPELINE_ENV)
}
